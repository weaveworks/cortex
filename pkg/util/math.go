package util

import (
	"github.com/cortexproject/cortex/pkg/util/mathutil"
)

// Max returns the maximum of two ints
func Max(a, b int) int {
	return mathutil.Max(a, b)
}

// Min returns the minimum of two ints
func Min(a, b int) int {
	return mathutil.Min(a, b)
}

// Max64 returns the maximum of two int64s
func Max64(a, b int64) int64 {
	return mathutil.Max64(a, b)
}

// Min64 returns the minimum of two int64s
func Min64(a, b int64) int64 {
	return mathutil.Min64(a, b)
}
