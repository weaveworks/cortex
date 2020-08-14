package stats

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
	"github.com/weaveworks/common/httpgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cortexproject/cortex/pkg/querier/frontend"
	"github.com/cortexproject/cortex/pkg/util/test"
)

func TestGrpcStats(t *testing.T) {
	reg := prometheus.NewRegistry()
	stats := NewStatsHandler(reg)

	serv := grpc.NewServer(grpc.StatsHandler(stats), grpc.MaxRecvMsgSize(10e6))
	defer serv.GracefulStop()

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go func() {
		require.NoError(t, serv.Serve(listener))
	}()

	grpc_health_v1.RegisterHealthServer(serv, health.NewServer())

	closed := false
	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)
	defer func() {
		if !closed {
			require.NoError(t, conn.Close())
		}
	}()

	hc := grpc_health_v1.NewHealthClient(conn)

	// First request (empty).
	resp, err := hc.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)

	// Second request, with large service name. This returns error, which doesn't count as "payload".
	_, err = hc.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
		Service: generateString(8 * 1024 * 1024),
	})
	require.EqualError(t, err, "rpc error: code = NotFound desc = unknown service")

	err = testutil.GatherAndCompare(reg, bytes.NewBufferString(`
			# HELP cortex_grpc_connected_clients Number of clients connected to gRPC server.
			# TYPE cortex_grpc_connected_clients gauge
			cortex_grpc_connected_clients 1

			# HELP cortex_grpc_method_errors_total Number of errors returned by method.
			# TYPE cortex_grpc_method_errors_total counter
			cortex_grpc_method_errors_total{method="/grpc.health.v1.Health/Check"} 1

			# HELP cortex_grpc_received_payload_size_bytes Size of received gRPC messages as seen on the wire (eg. compressed, signed, encrypted).
			# TYPE cortex_grpc_received_payload_size_bytes histogram
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="1.048576e+06"} 1
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="2.62144e+06"} 1
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="5.24288e+06"} 1
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="1.048576e+07"} 2
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="2.62144e+07"} 2
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="5.24288e+07"} 2
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="1.048576e+08"} 2
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="2.62144e+08"} 2
			cortex_grpc_received_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="+Inf"} 2
			cortex_grpc_received_payload_size_bytes_sum{method="/grpc.health.v1.Health/Check"} 8.388613e+06
			cortex_grpc_received_payload_size_bytes_count{method="/grpc.health.v1.Health/Check"} 2

			# HELP cortex_grpc_sent_payload_size_bytes Size of sent gRPC messages as seen on the wire (eg. compressed, signed, encrypted).
			# TYPE cortex_grpc_sent_payload_size_bytes histogram
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="1.048576e+06"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="2.62144e+06"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="5.24288e+06"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="1.048576e+07"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="2.62144e+07"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="5.24288e+07"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="1.048576e+08"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="2.62144e+08"} 1
			cortex_grpc_sent_payload_size_bytes_bucket{method="/grpc.health.v1.Health/Check",le="+Inf"} 1
			cortex_grpc_sent_payload_size_bytes_sum{method="/grpc.health.v1.Health/Check"} 7
			cortex_grpc_sent_payload_size_bytes_count{method="/grpc.health.v1.Health/Check"} 1
	`), "cortex_grpc_connected_clients", "cortex_grpc_received_payload_size_bytes", "cortex_grpc_sent_payload_size_bytes", "cortex_grpc_method_errors_total")
	require.NoError(t, err)

	closed = true
	require.NoError(t, conn.Close())

	// Give server little time to update connected clients metric.
	test.Poll(t, 1*time.Second, nil, func() interface{} {
		return testutil.GatherAndCompare(reg, bytes.NewBufferString(`
			# HELP cortex_grpc_connected_clients Number of clients connected to gRPC server.
			# TYPE cortex_grpc_connected_clients gauge
			cortex_grpc_connected_clients 0
		`), "cortex_grpc_connected_clients")
	})
}

func TestGrpcStatsStreaming(t *testing.T) {
	reg := prometheus.NewRegistry()
	stats := NewStatsHandler(reg)

	serv := grpc.NewServer(grpc.StatsHandler(stats), grpc.MaxSendMsgSize(10e6), grpc.MaxRecvMsgSize(10e6))
	defer serv.GracefulStop()

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go func() {
		require.NoError(t, serv.Serve(listener))
	}()

	frontend.RegisterFrontendServer(serv, &frontendServer{log: t.Log})

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10e6), grpc.MaxCallSendMsgSize(10e6)))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	fc := frontend.NewFrontendClient(conn)

	s, err := fc.Process(context.Background())
	require.NoError(t, err)

	for ix := 0; ix < 5; ix++ {
		req, err := s.Recv()
		require.NoError(t, err)

		msg := &frontend.ProcessResponse{HttpResponse: &httpgrpc.HTTPResponse{
			Code:    200,
			Headers: req.HttpRequest.Headers,
		}}
		t.Log("Client sending:", msg.Size())
		err = s.Send(msg)
		require.NoError(t, err)

		err = testutil.GatherAndCompare(reg, bytes.NewBufferString(`
			# HELP cortex_grpc_inflight_requests Number of inflight gRPC calls.
			# TYPE cortex_grpc_inflight_requests gauge
			cortex_grpc_inflight_requests{method="/frontend.Frontend/Process"} 1
		`), "cortex_grpc_inflight_requests")
		require.NoError(t, err)
	}
	require.NoError(t, s.CloseSend())

	// Wait until server notices.
	test.Poll(t, 1*time.Second, nil, func() interface{} {
		return testutil.GatherAndCompare(reg, bytes.NewBufferString(`
			# HELP cortex_grpc_inflight_requests Number of inflight gRPC calls.
			# TYPE cortex_grpc_inflight_requests gauge
			cortex_grpc_inflight_requests{method="/frontend.Frontend/Process"} 0
		`), "cortex_grpc_inflight_requests")
	})

	err = testutil.GatherAndCompare(reg, bytes.NewBufferString(`
			# HELP cortex_grpc_received_payload_size_bytes Size of received gRPC messages as seen on the wire (eg. compressed, signed, encrypted).
			# TYPE cortex_grpc_received_payload_size_bytes histogram
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="1.048576e+06"} 1
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="2.62144e+06"} 4
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="5.24288e+06"} 5
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="1.048576e+07"} 5
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="2.62144e+07"} 5
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="5.24288e+07"} 5
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="1.048576e+08"} 5
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="2.62144e+08"} 5
			cortex_grpc_received_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="+Inf"} 5
			cortex_grpc_received_payload_size_bytes_sum{method="/frontend.Frontend/Process"} 8.017448e+06
			cortex_grpc_received_payload_size_bytes_count{method="/frontend.Frontend/Process"} 5
			# HELP cortex_grpc_sent_payload_size_bytes Size of sent gRPC messages as seen on the wire (eg. compressed, signed, encrypted).
			# TYPE cortex_grpc_sent_payload_size_bytes histogram
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="1.048576e+06"} 0
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="2.62144e+06"} 2
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="5.24288e+06"} 4
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="1.048576e+07"} 6
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="2.62144e+07"} 6
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="5.24288e+07"} 6
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="1.048576e+08"} 6
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="2.62144e+08"} 6
			cortex_grpc_sent_payload_size_bytes_bucket{method="/frontend.Frontend/Process",le="+Inf"} 6
			cortex_grpc_sent_payload_size_bytes_sum{method="/frontend.Frontend/Process"} 2.2234511e+07
			cortex_grpc_sent_payload_size_bytes_count{method="/frontend.Frontend/Process"} 6
	`), "cortex_grpc_received_payload_size_bytes", "cortex_grpc_sent_payload_size_bytes")

	require.NoError(t, err)
}

type frontendServer struct {
	log func(args ...interface{})
}

func (f frontendServer) Process(server frontend.Frontend_ProcessServer) error {
	ix := 0
	for {
		ix++

		msg := &frontend.ProcessRequest{HttpRequest: &httpgrpc.HTTPRequest{
			Method: fmt.Sprintf("%d", ix),
			Url:    generateString(ix * 512 * 1024),
			Headers: []*httpgrpc.Header{
				{
					Key:    generateString(100 * ix),
					Values: []string{generateString(100 * ix), generateString(10000 * ix), generateString(ix * 512 * 1024)},
				},
			},
		}}

		f.log("Server sending:", msg.Size())
		err := server.Send(msg)

		if err != nil {
			return err
		}

		_, err = server.Recv()
		if err != nil {
			return err
		}
	}
}

func generateString(size int) string {
	// Use random bytes, to avoid compression.
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		// Should not happen.
		panic(err)
	}

	// To avoid invalid UTF-8 sequences (which protobuf complains about), we cleanup the data a bit.
	for ix, b := range buf {
		if b < ' ' {
			b += ' '
		}
		b = b & 0x7f
		buf[ix] = b
	}
	return string(buf)
}
