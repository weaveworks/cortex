package ingester

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/require"
	"github.com/weaveworks/common/middleware"
	"github.com/weaveworks/common/user"
	"google.golang.org/grpc"

	"github.com/cortexproject/cortex/pkg/chunk/encoding"
	"github.com/cortexproject/cortex/pkg/ingester/client"
	"github.com/cortexproject/cortex/pkg/ring"
	"github.com/cortexproject/cortex/pkg/storage/tsdb/backend/s3"
	"github.com/cortexproject/cortex/pkg/util/test"
	"github.com/cortexproject/cortex/pkg/util/validation"
)

func BenchmarkQueryStream(b *testing.B) {
	cfg := defaultIngesterTestConfig()
	clientCfg := defaultClientTestConfig()
	limits := defaultLimitsTestConfig()

	const (
		numSeries  = 1e6 // Put 1 million timeseries, each with 10 samples.
		numSamples = 10
		numCPUs    = 32
	)

	encoding.DefaultEncoding = encoding.Bigchunk
	limits.MaxLocalSeriesPerMetric = numSeries
	limits.MaxSeriesPerQuery = numSeries
	cfg.FlushCheckPeriod = 15 * time.Minute
	_, ing := newTestStore(b, cfg, clientCfg, limits)
	// defer ing.Shutdown()

	ctx := user.InjectOrgID(context.Background(), "1")
	instances := make([]string, numSeries/numCPUs)
	for i := 0; i < numSeries/numCPUs; i++ {
		instances[i] = fmt.Sprintf("node%04d", i)
	}
	cpus := make([]string, numCPUs)
	for i := 0; i < numCPUs; i++ {
		cpus[i] = fmt.Sprintf("cpu%02d", i)
	}

	for i := 0; i < numSeries; i++ {
		labels := labelPairs{
			{Name: model.MetricNameLabel, Value: "node_cpu"},
			{Name: "job", Value: "node_exporter"},
			{Name: "instance", Value: instances[i/numCPUs]},
			{Name: "cpu", Value: cpus[i%numCPUs]},
		}

		state, fp, series, err := ing.userStates.getOrCreateSeries(ctx, "1", labels, nil)
		require.NoError(b, err)

		for j := 0; j < numSamples; j++ {
			err = series.add(model.SamplePair{
				Value:     model.SampleValue(float64(j)),
				Timestamp: model.Time(int64(j)),
			})
			require.NoError(b, err)
		}

		state.fpLocker.Unlock(fp)
	}

	server := grpc.NewServer(grpc.StreamInterceptor(middleware.StreamServerUserHeaderInterceptor))
	defer server.GracefulStop()
	client.RegisterIngesterServer(server, ing)

	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(b, err)
	go server.Serve(l) //nolint:errcheck

	b.ResetTimer()
	for iter := 0; iter < b.N; iter++ {
		b.Run("QueryStream", func(b *testing.B) {
			c, err := client.MakeIngesterClient(l.Addr().String(), clientCfg)
			require.NoError(b, err)
			defer c.Close() //nolint:errcheck

			s, err := c.QueryStream(ctx, &client.QueryRequest{
				StartTimestampMs: 0,
				EndTimestampMs:   numSamples,
				Matchers: []*client.LabelMatcher{{
					Type:  client.EQUAL,
					Name:  model.MetricNameLabel,
					Value: "node_cpu",
				}},
			})
			require.NoError(b, err)

			count := 0
			for {
				resp, err := s.Recv()
				if err == io.EOF {
					break
				}
				require.NoError(b, err)
				count += len(resp.Timeseries)
			}
			require.Equal(b, count, int(numSeries))
		})
	}
}

func TestTSDBQueryStream(t *testing.T) {
	limits, err := validation.NewOverrides(defaultLimitsTestConfig(), nil)
	require.NoError(t, err)

	dir1, err := ioutil.TempDir("", "tsdb")
	require.NoError(t, err)

	// Start the first ingester, and get it into ACTIVE state.
	cfg1 := defaultIngesterTestConfig()
	cfg1.TSDBEnabled = true
	cfg1.TSDBConfig.Dir = dir1
	cfg1.TSDBConfig.S3 = s3.Config{
		Endpoint:        "dummy",
		BucketName:      "dummy",
		SecretAccessKey: "dummy",
		AccessKeyID:     "dummy",
	}
	cfg1.LifecyclerConfig.ID = "ingester1"
	cfg1.LifecyclerConfig.Addr = "ingester1"
	cfg1.LifecyclerConfig.JoinAfter = 0 * time.Second
	cfg1.MaxTransferRetries = 10
	ing, err := New(cfg1, defaultClientTestConfig(), limits, nil, nil)
	require.NoError(t, err)

	test.Poll(t, 100*time.Millisecond, ring.ACTIVE, func() interface{} {
		return ing.lifecycler.GetState()
	})

	// Now write a sample to this ingester
	const ts = 123000
	const val = 456
	var (
		l          = labels.Labels{{Name: labels.MetricName, Value: "foo"}}
		sampleData = []client.Sample{
			{
				TimestampMs: ts,
				Value:       val,
			},
		}
		expectedResponse = &client.QueryStreamResponse{
			Timeseries: []client.TimeSeries{
				{
					Labels: client.FromLabelsToLabelAdapters(l),
					Samples: []client.Sample{
						{
							Value:       val,
							TimestampMs: ts,
						},
					},
				},
			},
		}
	)

	serv := grpc.NewServer(grpc.StreamInterceptor(middleware.StreamServerUserHeaderInterceptor))
	defer serv.GracefulStop()
	client.RegisterIngesterServer(serv, ing)

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	go func() {
		err := serv.Serve(listener)
		require.NoError(t, err)
	}()

	// Write the metric
	ctx := user.InjectOrgID(context.Background(), userID)
	_, err = ing.Push(ctx, client.ToWriteRequest([]labels.Labels{l}, sampleData, client.API))
	require.NoError(t, err)

	// Stream the query
	c, err := client.MakeIngesterClient(listener.Addr().String(), defaultClientTestConfig())
	require.NoError(t, err)
	defer c.Close()

	s, err := c.QueryStream(ctx, &client.QueryRequest{
		StartTimestampMs: 0,
		EndTimestampMs:   200000,
		Matchers: []*client.LabelMatcher{{
			Type:  client.EQUAL,
			Name:  model.MetricNameLabel,
			Value: "foo",
		}},
	})
	require.NoError(t, err)

	count := 0
	var lastResp *client.QueryStreamResponse
	for {
		resp, err := s.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		count += len(resp.Timeseries)
		lastResp = resp
	}
	require.Equal(t, 1, count)
	require.Equal(t, expectedResponse, lastResp)
}
