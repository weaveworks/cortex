package gcp

import (
	"flag"
	"strings"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"

	"github.com/weaveworks/cortex/pkg/chunk"
)

const (
	columnFamily = "f"
	column       = "c"
	separator    = "\000"
)

// Config for a StorageClient
type Config struct {
	project  string
	instance string
}

// RegisterFlags adds the flags required to config this to the given FlagSet
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&cfg.project, "bigtable.project", "", "BigTable project ID.")
	f.StringVar(&cfg.instance, "bigtable.instance", "", "BigTable instance ID.")
}

// storageClient implements chunk.storageClient for GCP.
type storageClient struct {
	cfg       Config
	schemaCfg chunk.SchemaConfig
	client    *bigtable.Client
}

// NewStorageClient returns a new StorageClient.
func NewStorageClient(ctx context.Context, cfg Config, schemaCfg chunk.SchemaConfig) (chunk.StorageClient, error) {
	client, err := bigtable.NewClient(ctx, cfg.project, cfg.instance, instrumentation())
	if err != nil {
		return nil, err
	}
	return &storageClient{
		cfg:       cfg,
		schemaCfg: schemaCfg,
		client:    client,
	}, nil
}

func (s *storageClient) NewWriteBatch() chunk.WriteBatch {
	return bigtableWriteBatch{
		tables: map[string]map[string]*bigtable.Mutation{},
	}
}

type bigtableWriteBatch struct {
	tables map[string]map[string]*bigtable.Mutation
}

func (b bigtableWriteBatch) Add(tableName, hashValue string, rangeValue []byte, value []byte) {
	rows, ok := b.tables[tableName]
	if !ok {
		rows = map[string]*bigtable.Mutation{}
		b.tables[tableName] = rows
	}

	// TODO the hashValue should actually be hashed - but I have data written in
	// this format, so we need to do a proper migration.
	rowKey := hashValue + separator + string(rangeValue)
	mutation, ok := rows[rowKey]
	if !ok {
		mutation = bigtable.NewMutation()
		rows[rowKey] = mutation
	}

	mutation.Set(columnFamily, column, 0, value)
}

func (s *storageClient) BatchWrite(ctx context.Context, batch chunk.WriteBatch) error {
	bigtableBatch := batch.(bigtableWriteBatch)

	for tableName, rows := range bigtableBatch.tables {
		table := s.client.Open(tableName)
		rowKeys := make([]string, 0, len(rows))
		muts := make([]*bigtable.Mutation, 0, len(rows))
		for rowKey, mut := range rows {
			rowKeys = append(rowKeys, rowKey)
			muts = append(muts, mut)
		}

		errs, err := table.ApplyBulk(ctx, rowKeys, muts)
		if err != nil {
			return err
		}
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *storageClient) QueryPages(ctx context.Context, query chunk.IndexQuery, callback func(result chunk.ReadBatch, lastPage bool) (shouldContinue bool)) error {
	table := s.client.Open(query.TableName)

	var rowRange bigtable.RowRange
	if query.RangeValuePrefix != nil {
		rowRange = bigtable.PrefixRange(query.HashValue + separator + string(query.RangeValuePrefix))
	} else if query.RangeValueStart != nil {
		rowRange = bigtable.InfiniteRange(query.HashValue + separator + string(query.RangeValueStart))
	} else {
		rowRange = bigtable.PrefixRange(query.HashValue + separator)
	}

	return table.ReadRows(ctx, rowRange, func(r bigtable.Row) bool {
		// If RangeValueStart is set, BigTable doesn't know when to stop, as we're
		// reading "until the end of the row" in DynamoDB.  So we need to check
		// the prefix of the row is still correct.
		if query.RangeValueStart != nil {
			if !strings.HasPrefix(r.Key(), query.HashValue+separator) {
				return false
			}
		}

		return callback(bigtableReadBatch(r), false)
	}, bigtable.RowFilter(bigtable.FamilyFilter(columnFamily)))
}

// bigtableReadBatch represents a batch of rows read from BigTable.  As the
// bigtable interface gives us rows one-by-one, a batch always only contains
// a single row.
type bigtableReadBatch bigtable.Row

func (bigtableReadBatch) Len() int {
	return 1
}

func (b bigtableReadBatch) RangeValue(index int) []byte {
	if index != 0 {
		panic("index != 0")
	}
	// String before the first separator is the hashkey
	parts := strings.SplitN(bigtable.Row(b).Key(), separator, 2)
	return []byte(parts[1])
}

func (b bigtableReadBatch) Value(index int) []byte {
	if index != 0 {
		panic("index != 0")
	}
	cf, ok := b[columnFamily]
	if !ok || len(cf) != 1 {
		panic("bad response from bigtable")
	}
	return cf[0].Value
}

func (s *storageClient) PutChunks(ctx context.Context, chunks []chunk.Chunk) error {
	keys := map[string][]string{}
	muts := map[string][]*bigtable.Mutation{}

	for i := range chunks {
		// Encode the chunk first - checksum is calculated as a side effect.
		buf, err := chunks[i].Encode()
		if err != nil {
			return err
		}
		key := chunks[i].ExternalKey()
		tableName := s.schemaCfg.ChunkTables.TableFor(chunks[i].From)
		keys[tableName] = append(keys[tableName], key)

		mut := bigtable.NewMutation()
		mut.Set(columnFamily, column, 0, buf)
		muts[tableName] = append(muts[tableName], mut)
	}

	for tableName := range keys {
		table := s.client.Open(tableName)
		errs, err := table.ApplyBulk(ctx, keys[tableName], muts[tableName])
		if err != nil {
			return err
		}
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *storageClient) GetChunks(ctx context.Context, input []chunk.Chunk) ([]chunk.Chunk, error) {
	chunks := map[string][]chunk.Chunk{}
	keys := map[string]bigtable.RowList{}
	for _, chunk := range input {
		tableName := s.schemaCfg.ChunkTables.TableFor(chunk.From)
		keys[tableName] = append(keys[tableName], chunk.ExternalKey())
		chunks[tableName] = append(chunks[tableName], chunk)
	}

	output := make([]chunk.Chunk, 0, len(input))
	for tableName := range keys {
		var (
			i             = 0
			processingErr error
			table         = s.client.Open(tableName)
			chunks        = chunks[tableName]
		)
		// rows are returned in order
		if err := table.ReadRows(ctx, keys[tableName], func(row bigtable.Row) bool {
			chunk := chunks[i]
			i++
			processingErr = chunk.Decode(row[columnFamily][0].Value)
			if processingErr != nil {
				output = append(output, chunk)
			}
			return processingErr == nil
		}); err != nil {
			return nil, err
		}
		if processingErr != nil {
			return nil, processingErr
		}
	}
	return output, nil
}
