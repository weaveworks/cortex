package spanlogger

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSpanLogger_Log(t *testing.T) {
	span, ctx := New(context.Background(), "test", "bar")
	_ = span.Log("foo")
	newSpan := FromContext(ctx)
	require.Equal(t, span.Span, newSpan.Span)
	require.Error(t, newSpan.Log("bar", "err", errors.New("err"), "metric2", 2))
	noSpan := FromContext(context.Background())
	_ = noSpan.Log("foo")
	require.NoError(t, noSpan.Log("metric1", 1, "err", errors.New("err"), "metric2", 2))
}

func TestSpanLogger_CustomLogger(t *testing.T) {
	var logged [][]interface{}
	var logger funcLogger = func(keyvals ...interface{}) error {
		logged = append(logged, keyvals)
		return nil
	}
	span, ctx := NewWithLogger(context.Background(), logger, "test")
	_ = span.Log("msg", "original spanlogger")

	span = FromContextWithFallback(ctx, log.NewNopLogger())
	_ = span.Log("msg", "restored spanlogger")

	span = FromContextWithFallback(context.Background(), logger)
	_ = span.Log("msg", "fallback spanlogger")

	expect := [][]interface{}{
		{"method", "test", "msg", "original spanlogger"},
		{"msg", "restored spanlogger"},
		{"msg", "fallback spanlogger"},
	}
	require.Equal(t, expect, logged)
}

type funcLogger func(keyvals ...interface{}) error

func (f funcLogger) Log(keyvals ...interface{}) error {
	return f(keyvals...)
}
