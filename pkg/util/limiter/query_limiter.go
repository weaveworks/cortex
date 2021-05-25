package limiter

import (
	"context"
	"fmt"
	"sync"

	"github.com/prometheus/common/model"
	"go.uber.org/atomic"

	"github.com/cortexproject/cortex/pkg/cortexpb"
	"github.com/cortexproject/cortex/pkg/ingester/client"
	"github.com/cortexproject/cortex/pkg/util/validation"
)

type queryLimiterCtxKey struct{}

var (
	ctxKey              = &queryLimiterCtxKey{}
	errMaxSeriesHit     = "The query hit the max number of series limit (limit: %d)"
	errMaxChunkBytesHit = "The query hit the max number of chunk bytes limit (limit: %d)"
)

type QueryLimiter struct {
	uniqueSeriesMx sync.Mutex
	uniqueSeries   map[model.Fingerprint]struct{}

	chunkBytesCount *atomic.Int32

	maxSeriesPerQuery     int
	maxChunkBytesPerQuery int
}

// NewQueryLimiter makes a new per-query limiter. Each query limiter
// is configured using the `maxSeriesPerQuery` limit.
func NewQueryLimiter(maxSeriesPerQuery int, maxChunkBytesPerQuery int) *QueryLimiter {
	return &QueryLimiter{
		uniqueSeriesMx: sync.Mutex{},
		uniqueSeries:   map[model.Fingerprint]struct{}{},

		chunkBytesCount: atomic.NewInt32(0),

		maxSeriesPerQuery:     maxSeriesPerQuery,
		maxChunkBytesPerQuery: maxChunkBytesPerQuery,
	}
}

func AddQueryLimiterToContext(ctx context.Context, limiter *QueryLimiter) context.Context {
	return context.WithValue(ctx, ctxKey, limiter)
}

// QueryLimiterFromContextWithFallback returns a QueryLimiter from the current context.
// If there is not a QueryLimiter on the context it will return a new no-op limiter.
func QueryLimiterFromContextWithFallback(ctx context.Context) *QueryLimiter {
	ql, ok := ctx.Value(ctxKey).(*QueryLimiter)
	if !ok {
		// If there's no limiter return a new unlimited limiter as a fallback
		ql = NewQueryLimiter(0, 0)
	}
	return ql
}

// AddSeries adds the input series and returns an error if the limit is reached.
func (ql *QueryLimiter) AddSeries(seriesLabels []cortexpb.LabelAdapter) error {
	// If the max series is unlimited just return without managing map
	if ql.maxSeriesPerQuery == 0 {
		return nil
	}
	fingerprint := client.FastFingerprint(seriesLabels)

	ql.uniqueSeriesMx.Lock()
	defer ql.uniqueSeriesMx.Unlock()

	ql.uniqueSeries[fingerprint] = struct{}{}
	if len(ql.uniqueSeries) > ql.maxSeriesPerQuery {
		// Format error with max limit
		return validation.LimitError(fmt.Sprintf(errMaxSeriesHit, ql.maxSeriesPerQuery))
	}
	return nil
}

// uniqueSeriesCount returns the count of unique series seen by this query limiter.
func (ql *QueryLimiter) uniqueSeriesCount() int {
	ql.uniqueSeriesMx.Lock()
	defer ql.uniqueSeriesMx.Unlock()
	return len(ql.uniqueSeries)
}

func (ql *QueryLimiter) AddChunkBytes(bytes int) error {
	if ql.maxChunkBytesPerQuery == 0 {
		return nil
	}
	if ql.chunkBytesCount.Add(int32(bytes)) > int32(ql.maxChunkBytesPerQuery) {
		return validation.LimitError(fmt.Sprintf(errMaxChunkBytesHit, ql.maxChunkBytesPerQuery))
	}
	return nil
}
