package chunk

import (
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/local/chunk"
	"github.com/stretchr/testify/require"
)

const userID = "userID"

func dummyChunk() Chunk {
	return dummyChunkFor(model.Metric{
		model.MetricNameLabel: "foo",
		"bar":  "baz",
		"toms": "code",
	})
}

func dummyChunkFor(metric model.Metric) Chunk {
	now := model.Now()
	cs, _ := chunk.New().Add(model.SamplePair{Timestamp: now, Value: 0})
	chunk := NewChunk(
		userID,
		metric.Fingerprint(),
		metric,
		cs[0],
		now.Add(-time.Hour),
		now,
	)
	return chunk
}

func TestChunkCodec(t *testing.T) {
	for _, c := range []struct {
		chunk Chunk
		err   error
		f     func(*Chunk, []byte)
	}{
		// Basic round trip
		{chunk: dummyChunk()},

		// Checksum should fail
		{
			chunk: dummyChunk(),
			err:   ErrInvalidChecksum,
			f:     func(_ *Chunk, buf []byte) { buf[4] += 1 },
		},

		// Checksum should fail
		{
			chunk: dummyChunk(),
			err:   ErrInvalidChecksum,
			f:     func(c *Chunk, _ []byte) { c.Checksum = 123 },
		},

		// Metadata test should fail
		{
			chunk: dummyChunk(),
			err:   ErrWrongMetadata,
			f:     func(c *Chunk, _ []byte) { c.Fingerprint += 1 },
		},
	} {
		buf, err := c.chunk.encode()
		require.NoError(t, err)

		have, err := parseExternalKey("", c.chunk.externalKey())
		require.NoError(t, err)

		if c.f != nil {
			c.f(&have, buf)
		}

		err = have.decode(buf)
		require.Equal(t, err, c.err)

		if c.err == nil {
			require.Equal(t, have, c.chunk)
		}
	}
}

func TestParseExternalKey(t *testing.T) {
	for _, c := range []struct {
		key   string
		chunk Chunk
		err   error
	}{
		{key: "2:1484661279394:1484664879394", chunk: Chunk{
			UserID:      userID,
			Fingerprint: model.Fingerprint(2),
			From:        model.Time(1484661279394),
			Through:     model.Time(1484664879394),
		}},

		{key: "1/2:270d8f00:270d8f00:f84c5745", chunk: Chunk{
			UserID:      "1",
			Fingerprint: model.Fingerprint(2),
			From:        model.Time(655200000),
			Through:     model.Time(655200000),
			ChecksumSet: true,
			Checksum:    4165752645,
		}},
	} {
		chunk, err := parseExternalKey(userID, c.key)
		require.Equal(t, c.err, err)
		require.Equal(t, c.chunk, chunk)
	}
}
