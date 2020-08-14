package stats

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/stats"
)

func NewStatsHandler(r prometheus.Registerer) stats.Handler {
	const MiB = 1024 * 1024
	messageSizeBuckets := []float64{1 * MiB, 2.5 * MiB, 5 * MiB, 10 * MiB, 25 * MiB, 50 * MiB, 100 * MiB, 250 * MiB}

	return &grpcStatsHandler{
		connectedClients: promauto.With(r).NewGauge(prometheus.GaugeOpts{
			Name: "cortex_grpc_connected_clients",
			Help: "Number of clients connected to gRPC server.",
		}),

		inflightRPC: promauto.With(r).NewGaugeVec(prometheus.GaugeOpts{
			Name: "cortex_grpc_inflight_requests",
			Help: "Number of inflight gRPC calls.",
		}, []string{"method"}),

		methodErrors: promauto.With(r).NewCounterVec(prometheus.CounterOpts{
			Name: "cortex_grpc_method_errors_total",
			Help: "Number of errors returned by method.",
		}, []string{"method"}),

		receivedPayloadSize: promauto.With(r).NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cortex_grpc_received_payload_size_bytes",
			Help:    "Size of received gRPC messages as seen on the wire (eg. compressed, signed, encrypted).",
			Buckets: messageSizeBuckets,
		}, []string{"method"}),

		sentPayloadSize: promauto.With(r).NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cortex_grpc_sent_payload_size_bytes",
			Help:    "Size of sent gRPC messages as seen on the wire (eg. compressed, signed, encrypted).",
			Buckets: messageSizeBuckets,
		}, []string{"method"}),
	}
}

type grpcStatsHandler struct {
	connectedClients    prometheus.Gauge
	inflightRPC         *prometheus.GaugeVec
	receivedPayloadSize *prometheus.HistogramVec
	sentPayloadSize     *prometheus.HistogramVec
	methodErrors        *prometheus.CounterVec
}

// Custom type to hide it from other packages.
type contextKey int

const (
	contextKeyMethodName contextKey = 1
)

func (g *grpcStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return context.WithValue(ctx, contextKeyMethodName, info.FullMethodName)
}

func (g *grpcStatsHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	// We use full method name from context, because not all RPCStats structs have it.
	fullMethodName, ok := ctx.Value(contextKeyMethodName).(string)
	if !ok {
		return
	}

	switch s := rpcStats.(type) {
	case *stats.Begin:
		g.inflightRPC.WithLabelValues(fullMethodName).Inc()
	case *stats.End:
		g.inflightRPC.WithLabelValues(fullMethodName).Dec()
		if s.Error != nil {
			g.methodErrors.WithLabelValues(fullMethodName).Inc()
		}

	case *stats.InHeader:
		// Ignored. Cortex doesn't use headers. Furthermore WireLength seems to be incorrect for large headers -- it uses
		// length of last frame (16K) even for headers in megabytes.
	case *stats.InPayload:
		g.receivedPayloadSize.WithLabelValues(fullMethodName).Observe(float64(s.WireLength))
	case *stats.InTrailer:
		// Ignored. Cortex doesn't use trailers.

	case *stats.OutHeader:
		// Ignored. Cortex doesn't send headers, and since OutHeader doesn't have WireLength, we could only estimate it.
	case *stats.OutPayload:
		g.sentPayloadSize.WithLabelValues(fullMethodName).Observe(float64(s.WireLength))
	case *stats.OutTrailer:
		// Ignored, Cortex doesn't use trailers. OutTrailer doesn't have valid WireLength (there is deperecated field, always set to 0).
	}
}

func (g *grpcStatsHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (g *grpcStatsHandler) HandleConn(_ context.Context, connStats stats.ConnStats) {
	switch connStats.(type) {
	case *stats.ConnBegin:
		g.connectedClients.Inc()

	case *stats.ConnEnd:
		g.connectedClients.Dec()
	}
}
