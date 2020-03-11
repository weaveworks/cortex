// +build integration

package main

import (
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cortexproject/cortex/integration/e2e"
	e2edb "github.com/cortexproject/cortex/integration/e2e/db"
	"github.com/cortexproject/cortex/integration/e2ecortex"
)

func TestAllIndexStores(t *testing.T) {
	s, err := e2e.NewScenario(networkName)
	require.NoError(t, err)
	defer s.Close()

	// Start dependencies.
	dynamo := e2edb.NewDynamoDB()
	bigtable := e2edb.NewBigtable()

	stores := []string{"aws-dynamo", "bigtable"}
	perStoreDuration := 14 * 24 * time.Hour

	consul := e2edb.NewConsul()
	require.NoError(t, s.StartAndWaitReady(dynamo, bigtable, consul))

	// lets build config for each type of Index Store.
	now := time.Now()
	oldestStoreStartTime := now.Add(time.Duration(-len(stores)) * perStoreDuration)

	storeConfigs := make([]storeConfig, len(stores))
	for i, store := range stores {
		storeConfigs[i] = storeConfig{From: oldestStoreStartTime.Add(time.Duration(i) * perStoreDuration).Format("2006-01-02"), IndexStore: store}
	}

	// bigtable client needs to set an environment variable when connecting to an emulator
	bigtableFlag := map[string]string{"BIGTABLE_EMULATOR_HOST": bigtable.NetworkHTTPEndpoint()}

	// here we are starting and stopping table manager for each index store
	// this is a workaround to make table manager create tables for each config since it considers only latest schema config while creating tables
	for i := range storeConfigs {
		require.NoError(t, writeFileToSharedDir(s, cortexSchemaConfigFile, []byte(buildSchemaConfigWith(storeConfigs[i:i+1]))))

		tableManager := e2ecortex.NewTableManager("table-manager", mergeFlags(ChunksStorageFlags, map[string]string{
			"-table-manager.retention-period": "2520h", // setting retention high enough
		}), "", bigtableFlag)
		require.NoError(t, s.StartAndWaitReady(tableManager))

		// Wait until the first table-manager sync has completed, so that we're
		// sure the tables have been created.
		require.NoError(t, tableManager.WaitSumMetrics(e2e.Greater(0), "cortex_dynamo_sync_tables_seconds"))
		require.NoError(t, s.Stop(tableManager))
	}

	// Start rest of the Cortex components.
	require.NoError(t, writeFileToSharedDir(s, cortexSchemaConfigFile, []byte(buildSchemaConfigWith(storeConfigs))))

	ingester := e2ecortex.NewIngester("ingester", consul.NetworkHTTPEndpoint(), ChunksStorageFlags, "", bigtableFlag)
	distributor := e2ecortex.NewDistributor("distributor", consul.NetworkHTTPEndpoint(), ChunksStorageFlags, "")
	querier := e2ecortex.NewQuerier("querier", consul.NetworkHTTPEndpoint(), ChunksStorageFlags, "", bigtableFlag)

	require.NoError(t, s.StartAndWaitReady(distributor, ingester, querier))

	// Wait until both the distributor and querier have updated the ring.
	require.NoError(t, distributor.WaitSumMetrics(e2e.Equals(512), "cortex_ring_tokens_total"))
	require.NoError(t, querier.WaitSumMetrics(e2e.Equals(512), "cortex_ring_tokens_total"))

	// Push and Query some series to Cortex for each day starting from oldest start time from configs until now so that we test all the Index Stores
	for ts := oldestStoreStartTime; ts.Before(now); ts = ts.Add(24 * time.Hour) {
		series, expectedVector := generateSeries("series_1", ts)

		c, err := e2ecortex.NewClient(distributor.HTTPEndpoint(), "", "", "user-1")
		require.NoError(t, err)

		res, err := c.Push(series)
		require.NoError(t, err)
		require.Equal(t, 200, res.StatusCode)

		// Query the series both from the querier and query-frontend (to hit the read path).
		c, err = e2ecortex.NewClient("", querier.HTTPEndpoint(), "", "user-1")
		require.NoError(t, err)

		result, err := c.Query("series_1", ts)
		require.NoError(t, err)
		require.Equal(t, model.ValVector, result.Type())
		assert.Equal(t, expectedVector, result.(model.Vector))
	}
}
