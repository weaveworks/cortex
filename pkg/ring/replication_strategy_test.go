package ring

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRingReplicationStrategy(t *testing.T) {
	for i, tc := range []struct {
		replifcationFactor, liveIngesters, deadIngesters int
		expectedMaxFailure                               int
		expectedError                                    string
	}{
		// Ensure it works for a single ingester, for local testing.
		{
			replifcationFactor: 1,
			liveIngesters:      1,
			expectedMaxFailure: 0,
		},

		{
			replifcationFactor: 1,
			deadIngesters:      1,
			expectedError:      "at least 1 live replicas required, could only find 0",
		},

		// Ensure it works for RF=3 and 2 ingesters.
		{
			replifcationFactor: 3,
			liveIngesters:      2,
			expectedMaxFailure: 0,
		},

		// Ensure it works for the default production config.
		{
			replifcationFactor: 3,
			liveIngesters:      3,
			expectedMaxFailure: 1,
		},

		{
			replifcationFactor: 3,
			liveIngesters:      2,
			deadIngesters:      1,
			expectedMaxFailure: 0,
		},

		{
			replifcationFactor: 3,
			liveIngesters:      1,
			deadIngesters:      2,
			expectedError:      "at least 2 live replicas required, could only find 1",
		},

		// Ensure it works when adding / removing nodes.

		// A node is joining or leaving, replica set expands.
		{
			replifcationFactor: 3,
			liveIngesters:      4,
			expectedMaxFailure: 1,
		},

		{
			replifcationFactor: 3,
			liveIngesters:      3,
			deadIngesters:      1,
			expectedMaxFailure: 0,
		},

		{
			replifcationFactor: 3,
			liveIngesters:      2,
			deadIngesters:      2,
			expectedError:      "at least 3 live replicas required, could only find 2",
		},
	} {
		ingesters := []IngesterDesc{}
		for i := 0; i < tc.liveIngesters; i++ {
			ingesters = append(ingesters, IngesterDesc{
				Timestamp: time.Now().Unix(),
			})
		}
		for i := 0; i < tc.deadIngesters; i++ {
			ingesters = append(ingesters, IngesterDesc{})
		}

		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			strategy := NewDefaultReplicationStrategy()
			liveIngesters, maxFailure, err := strategy.Filter(ingesters, Read, tc.RF, 100*time.Second, false)
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.liveIngesters, len(liveIngesters))
				assert.Equal(t, tc.expectedMaxFailure, maxFailure)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestIgnoreUnhealthyInstancesReplicationStrategy(t *testing.T) {
	for _, tc := range []struct {
		name                         string
		liveIngesters, deadIngesters int
		expectedMaxFailure           int
		expectedError                string
	}{
		{
			name:               "with at least 1 healthy instance",
			liveIngesters:      1,
			expectedMaxFailure: 0,
		},
		{
			name:               "with more healthy instances than unhealthy",
			deadIngesters:      1,
			liveIngesters:      2,
			expectedMaxFailure: 1,
		},
		{
			name:               "with more unhealthy instances than healthy",
			deadIngesters:      1,
			liveIngesters:      2,
			expectedMaxFailure: 1,
		},
		{
			name:               "with equal number of healthy and unhealthy instances",
			deadIngesters:      2,
			liveIngesters:      2,
			expectedMaxFailure: 1,
		},
		{
			name:               "with no healthy instances",
			liveIngesters:      0,
			deadIngesters:      3,
			expectedMaxFailure: 0,
			expectedError:      "at least 1 healthy replica required, could only find 0",
		},
	} {
		ingesters := []IngesterDesc{}
		for i := 0; i < tc.liveIngesters; i++ {
			ingesters = append(ingesters, IngesterDesc{
				Timestamp: time.Now().Unix(),
			})
		}
		for i := 0; i < tc.deadIngesters; i++ {
			ingesters = append(ingesters, IngesterDesc{})
		}

		t.Run(tc.name, func(t *testing.T) {
			strategy := NewIgnoreUnhealthyInstancesReplicationStrategy()
			liveIngesters, maxFailure, err := strategy.Filter(ingesters, Read, 3, 100*time.Second, false)
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.liveIngesters, len(liveIngesters))
				assert.Equal(t, tc.expectedMaxFailure, maxFailure)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
