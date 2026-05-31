package telemetry

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var once sync.Once

// Init registers the custom database collector with a 60-second cache TTL.
func Init() {
	once.Do(func() {
		collector := NewDatabaseCollector(60 * time.Second)
		prometheus.MustRegister(collector)
	})
}
