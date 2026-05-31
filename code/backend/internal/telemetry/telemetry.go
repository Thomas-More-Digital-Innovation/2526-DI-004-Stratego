package telemetry

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Init registers the custom database collector with a 60-second cache TTL.
func Init() {
	collector := NewDatabaseCollector(60 * time.Second)
	prometheus.MustRegister(collector)
}
