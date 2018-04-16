package testutils

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
	promchunk "github.com/weaveworks/cortex/pkg/prom1/storage/local/chunk"

	"github.com/weaveworks/cortex/pkg/chunk"
)

const (
	userID = "userID"
)

// Fixture type for per-backend testing.
type Fixture interface {
	Name() string
	Clients() (chunk.StorageClient, chunk.TableClient, chunk.SchemaConfig, error)
	Teardown() error
}

func Setup(fixture Fixture, tableName string) (chunk.StorageClient, error) {
	storageClient, tableClient, schemaConfig, err := fixture.Clients()
	if err != nil {
		return nil, err
	}

	tableManager, err := chunk.NewTableManager(schemaConfig, 12*time.Hour, tableClient)
	if err != nil {
		return nil, err
	}

	err = tableManager.SyncTables(context.Background())
	if err != nil {
		return nil, err
	}

	err = tableClient.CreateTable(context.Background(), chunk.TableDesc{
		Name: tableName,
	})
	return storageClient, err
}

func CreateChunks(startIndex, batchSize int) ([]string, []chunk.Chunk, error) {
	keys := []string{}
	chunks := []chunk.Chunk{}
	for j := 0; j < batchSize; j++ {
		chunk := dummyChunkFor(model.Now(), model.Metric{
			model.MetricNameLabel: "foo",
			"index":               model.LabelValue(strconv.Itoa(startIndex*batchSize + j)),
		})
		chunks = append(chunks, chunk)
		_, err := chunk.Encode() // Need to encode it, side effect calculates crc
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, chunk.ExternalKey())
	}
	return keys, chunks, nil
}

func dummyChunk(now model.Time) chunk.Chunk {
	return dummyChunkFor(now, model.Metric{
		model.MetricNameLabel: "foo",
		"bar":  "baz",
		"toms": "code",
	})
}

func dummyChunkFor(now model.Time, metric model.Metric) chunk.Chunk {
	cs, _ := promchunk.New().Add(model.SamplePair{Timestamp: now, Value: 0})
	chunk := chunk.NewChunk(
		userID,
		metric.Fingerprint(),
		metric,
		cs[0],
		now.Add(-time.Hour),
		now,
	)
	// Force checksum calculation.
	_, err := chunk.Encode()
	if err != nil {
		panic(err)
	}
	return chunk
}
