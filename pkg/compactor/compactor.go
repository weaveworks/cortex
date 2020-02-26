package compactor

import (
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"path"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/pstibrany/services"
	"github.com/thanos-io/thanos/pkg/block"
	"github.com/thanos-io/thanos/pkg/compact"
	"github.com/thanos-io/thanos/pkg/compact/downsample"
	"github.com/thanos-io/thanos/pkg/objstore"

	"github.com/cortexproject/cortex/pkg/ring"
	cortex_tsdb "github.com/cortexproject/cortex/pkg/storage/tsdb"
	"github.com/cortexproject/cortex/pkg/util"
)

// Config holds the Compactor config.
type Config struct {
	BlockRanges          cortex_tsdb.DurationList `yaml:"block_ranges"`
	BlockSyncConcurrency int                      `yaml:"block_sync_concurrency"`
	MetaSyncConcurrency  int                      `yaml:"meta_sync_concurrency"`
	ConsistencyDelay     time.Duration            `yaml:"consistency_delay"`
	DataDir              string                   `yaml:"data_dir"`
	CompactionInterval   time.Duration            `yaml:"compaction_interval"`
	CompactionRetries    int                      `yaml:"compaction_retries"`

	// Compactors sharding.
	ShardingEnabled bool       `yaml:"sharding_enabled"`
	ShardingRing    RingConfig `yaml:"sharding_ring"`

	// No need to add options to customize the retry backoff,
	// given the defaults should be fine, but allow to override
	// it in tests.
	retryMinBackoff time.Duration `yaml:"-"`
	retryMaxBackoff time.Duration `yaml:"-"`
}

// RegisterFlags registers the Compactor flags.
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.ShardingRing.RegisterFlags(f)

	cfg.BlockRanges = cortex_tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}
	cfg.retryMinBackoff = 10 * time.Second
	cfg.retryMaxBackoff = time.Minute

	f.Var(&cfg.BlockRanges, "compactor.block-ranges", "List of compaction time ranges.")
	f.DurationVar(&cfg.ConsistencyDelay, "compactor.consistency-delay", 30*time.Minute, fmt.Sprintf("Minimum age of fresh (non-compacted) blocks before they are being processed. Malformed blocks older than the maximum of consistency-delay and %s will be removed.", compact.PartialUploadThresholdAge))
	f.IntVar(&cfg.BlockSyncConcurrency, "compactor.block-sync-concurrency", 20, "Number of Go routines to use when syncing block index and chunks files from the long term storage.")
	f.IntVar(&cfg.MetaSyncConcurrency, "compactor.meta-sync-concurrency", 20, "Number of Go routines to use when syncing block meta files from the long term storage.")
	f.StringVar(&cfg.DataDir, "compactor.data-dir", "./data", "Data directory in which to cache blocks and process compactions")
	f.DurationVar(&cfg.CompactionInterval, "compactor.compaction-interval", time.Hour, "The frequency at which the compaction runs")
	f.IntVar(&cfg.CompactionRetries, "compactor.compaction-retries", 3, "How many times to retry a failed compaction during a single compaction interval")
	f.BoolVar(&cfg.ShardingEnabled, "compactor.sharding-enabled", false, "Shard tenants across multiple compactor instances. Sharding is required if you run multiple compactor instances, in order to coordinate compactions and avoid race conditions leading to the same tenant blocks simultaneously compacted by different instances.")
}

// Compactor is a multi-tenant TSDB blocks compactor based on Thanos.
type Compactor struct {
	services.BasicService

	compactorCfg Config
	storageCfg   cortex_tsdb.Config
	logger       log.Logger

	// function that creates bucket client and TSDB compactor using the context.
	// Useful for injecting mock objects from tests.
	createBucketClientAndTsdbCompactor func(ctx context.Context) (objstore.Bucket, tsdb.Compactor, error)

	// Underlying compactor used to compact TSDB blocks.
	tsdbCompactor tsdb.Compactor

	// Client used to run operations on the bucket storing blocks.
	bucketClient objstore.Bucket

	// Ring used for sharding compactions.
	ringLifecycler *ring.Lifecycler
	ring           *ring.Ring

	// Manager sub-services (ring, lifecycler)
	subservices *services.Manager

	// Metrics.
	compactionRunsStarted   prometheus.Counter
	compactionRunsCompleted prometheus.Counter
	compactionRunsFailed    prometheus.Counter

	// TSDB syncer metrics
	syncerMetrics *syncerMetrics
}

// NewCompactor makes a new Compactor.
func NewCompactor(compactorCfg Config, storageCfg cortex_tsdb.Config, logger log.Logger, registerer prometheus.Registerer) (*Compactor, error) {
	createObjectsFn := func(ctx context.Context) (objstore.Bucket, tsdb.Compactor, error) {
		bucketClient, err := cortex_tsdb.NewBucketClient(ctx, storageCfg, "compactor", logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create the bucket client")
		}

		if registerer != nil {
			bucketClient = objstore.BucketWithMetrics( /* bucket label value */ "", bucketClient, prometheus.WrapRegistererWithPrefix("cortex_compactor_", registerer))
		}

		compactor, err := tsdb.NewLeveledCompactor(ctx, registerer, logger, compactorCfg.BlockRanges.ToMilliseconds(), downsample.NewPool())
		return bucketClient, compactor, err
	}

	cortexCompactor, err := newCompactor(compactorCfg, storageCfg, logger, registerer, createObjectsFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Cortex blocks compactor")
	}

	return cortexCompactor, nil
}

func newCompactor(
	compactorCfg Config,
	storageCfg cortex_tsdb.Config,
	logger log.Logger,
	registerer prometheus.Registerer,
	createTsdbCompactor func(ctx context.Context) (objstore.Bucket, tsdb.Compactor, error),
) (*Compactor, error) {
	c := &Compactor{
		compactorCfg:                       compactorCfg,
		storageCfg:                         storageCfg,
		logger:                             logger,
		createBucketClientAndTsdbCompactor: createTsdbCompactor,

		compactionRunsStarted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cortex_compactor_runs_started_total",
			Help: "Total number of compaction runs started.",
		}),
		compactionRunsCompleted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cortex_compactor_runs_completed_total",
			Help: "Total number of compaction runs successfully completed.",
		}),
		compactionRunsFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cortex_compactor_runs_failed_total",
			Help: "Total number of compaction runs failed.",
		}),
	}

	// Register metrics.
	if registerer != nil {
		registerer.MustRegister(c.compactionRunsStarted, c.compactionRunsCompleted, c.compactionRunsFailed)
		c.syncerMetrics = newSyncerMetrics(registerer)
	}

	services.InitBasicService(&c.BasicService, c.starting, c.running, c.stopping)

	return c, nil
}

// Start the compactor.
func (c *Compactor) starting(ctx context.Context) error {
	// Initialize the compactors ring if sharding is enabled.
	if c.compactorCfg.ShardingEnabled {
		lifecyclerCfg := c.compactorCfg.ShardingRing.ToLifecyclerConfig()
		lifecycler, err := ring.NewLifecycler(lifecyclerCfg, ring.NewNoopFlushTransferer(), "compactor", ring.CompactorRingKey, false)
		if err != nil {
			return errors.Wrap(err, "unable to initialize compactor ring lifecycler")
		}

		c.ringLifecycler = lifecycler

		ring, err := ring.New(lifecyclerCfg.RingConfig, "compactor", ring.CompactorRingKey)
		if err != nil {
			return errors.Wrap(err, "unable to initialize compactor ring")
		}

		c.ring = ring

		c.subservices, err = services.NewManager(c.ringLifecycler, c.ring)
		if err == nil {
			err = c.subservices.StartAsync(ctx)
			if err == nil {
				err = c.subservices.AwaitHealthy(ctx)
			}
		}

		if err != nil {
			return errors.Wrap(err, "unable to initialize service manager")
		}
	}

	var err error
	c.bucketClient, c.tsdbCompactor, err = c.createBucketClientAndTsdbCompactor(ctx)
	if err != nil && c.subservices != nil {
		c.subservices.StopAsync()
	}

	return errors.Wrap(err, "failed to initialize compactor objects")
}

func (c *Compactor) stopping() error {
	if c.subservices != nil {
		c.subservices.StopAsync()
		_ = c.subservices.AwaitStopped(context.Background())
	}

	return nil
}

func (c *Compactor) running(ctx context.Context) error {
	// If sharding is enabled we should wait until this instance is
	// ACTIVE within the ring.
	if c.compactorCfg.ShardingEnabled {
		level.Info(c.logger).Log("msg", "waiting until compactor is ACTIVE in the ring")
		if err := c.waitRingActive(ctx); err != nil {
			return err
		}
		level.Info(c.logger).Log("msg", "compactor is ACTIVE in the ring")
	}

	// Run an initial compaction before starting the interval.
	c.compactUsersWithRetries(ctx)

	ticker := time.NewTicker(c.compactorCfg.CompactionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.compactUsersWithRetries(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Compactor) compactUsersWithRetries(ctx context.Context) {
	retries := util.NewBackoff(ctx, util.BackoffConfig{
		MinBackoff: c.compactorCfg.retryMinBackoff,
		MaxBackoff: c.compactorCfg.retryMaxBackoff,
		MaxRetries: c.compactorCfg.CompactionRetries,
	})

	c.compactionRunsStarted.Inc()

	for retries.Ongoing() {
		if success := c.compactUsers(ctx); success {
			c.compactionRunsCompleted.Inc()
			return
		}

		retries.Wait()
	}

	c.compactionRunsFailed.Inc()
}

func (c *Compactor) compactUsers(ctx context.Context) bool {
	level.Info(c.logger).Log("msg", "discovering users from bucket")
	users, err := c.discoverUsers(ctx)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to discover users from bucket", "err", err)
		return false
	}
	level.Info(c.logger).Log("msg", "discovered users from bucket", "users", len(users))

	for _, userID := range users {
		// Ensure the context has not been canceled (ie. compactor shutdown has been triggered).
		if ctx.Err() != nil {
			level.Info(c.logger).Log("msg", "interrupting compaction of user blocks", "err", err)
			return false
		}

		// If sharding is enabled, ensure the user ID belongs to our shard.
		if c.compactorCfg.ShardingEnabled {
			if owned, err := c.ownUser(userID); err != nil {
				level.Warn(c.logger).Log("msg", "unable to check if user is owned by this shard", "user", userID, "err", err)
				continue
			} else if !owned {
				level.Debug(c.logger).Log("msg", "skipping user because not owned by this shard", "user", userID)
				continue
			}
		}

		level.Info(c.logger).Log("msg", "starting compaction of user blocks", "user", userID)

		if err = c.compactUser(ctx, userID); err != nil {
			level.Error(c.logger).Log("msg", "failed to compact user blocks", "user", userID, "err", err)
			continue
		}

		level.Info(c.logger).Log("msg", "successfully compacted user blocks", "user", userID)
	}

	return true
}

func (c *Compactor) compactUser(ctx context.Context, userID string) error {
	bucket := cortex_tsdb.NewUserBucketClient(userID, c.bucketClient)

	reg := prometheus.NewRegistry()
	defer c.syncerMetrics.gatherThanosSyncerMetrics(reg)

	ulogger := util.WithUserID(userID, c.logger)

	fetcher, err := block.NewMetaFetcher(
		ulogger,
		c.compactorCfg.MetaSyncConcurrency,
		bucket,
		// The fetcher stores cached metas in the "meta-syncer/" sub directory,
		// but we prefix it with "meta-" in order to guarantee no clashing with
		// the directory used by the Thanos Syncer, whatever is the user ID.
		path.Join(c.compactorCfg.DataDir, "meta-"+userID),
		reg,
		// No filters
	)
	if err != nil {
		return err
	}

	syncer, err := compact.NewSyncer(
		ulogger,
		reg,
		bucket,
		fetcher,
		c.compactorCfg.BlockSyncConcurrency,
		false, // Do not accept malformed indexes
		true,  // Enable vertical compaction
	)
	if err != nil {
		return errors.Wrap(err, "failed to create syncer")
	}

	compactor, err := compact.NewBucketCompactor(
		ulogger,
		syncer,
		c.tsdbCompactor,
		path.Join(c.compactorCfg.DataDir, "compact"),
		bucket,
		// No compaction concurrency. Due to how Cortex works we don't
		// expect to have multiple block groups per tenant, so setting
		// a value higher than 1 would be useless.
		1,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create bucket compactor")
	}

	return compactor.Compact(ctx)
}

func (c *Compactor) discoverUsers(ctx context.Context) ([]string, error) {
	var users []string

	err := c.bucketClient.Iter(ctx, "", func(entry string) error {
		users = append(users, strings.TrimSuffix(entry, "/"))
		return nil
	})

	return users, err
}

func (c *Compactor) ownUser(userID string) (bool, error) {
	// Hash the user ID.
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(userID))
	userHash := hasher.Sum32()

	// Check whether this compactor instance owns the user.
	rs, err := c.ring.Get(userHash, ring.Read, []ring.IngesterDesc{})
	if err != nil {
		return false, err
	}

	if len(rs.Ingesters) != 1 {
		return false, fmt.Errorf("unexpected number of compactors in the shard (expected 1, got %d)", len(rs.Ingesters))
	}

	return rs.Ingesters[0].Addr == c.ringLifecycler.Addr, nil
}

func (c *Compactor) waitRingActive(ctx context.Context) error {
	for {
		// Check if the ingester is ACTIVE in the ring and our ring client
		// has detected it.
		if rs, err := c.ring.GetAll(); err == nil {
			for _, i := range rs.Ingesters {
				if i.GetAddr() == c.ringLifecycler.Addr && i.GetState() == ring.ACTIVE {
					return nil
				}
			}
		}

		select {
		case <-time.After(time.Second):
			// Nothing to do
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
