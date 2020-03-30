package util

import (
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var experimentalModulesInUse = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "cortex",
		Name:      "experimental_features_in_use_total",
		Help:      "The number of experimental features in use.",
	},
)

// WarnExperimentalUse logs a warning and increments the experimental features metric.
func WarnExperimentalUse(module string) {
	level.Warn(Logger).Log("msg", "experimental feature in use", "module", module)
	experimentalModulesInUse.Inc()
}
