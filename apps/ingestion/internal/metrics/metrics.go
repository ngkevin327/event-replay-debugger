package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	batchLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "ingest_batch_latency_seconds",
		Help:    "Latency of ingest batch handling",
		Buckets: prometheus.DefBuckets,
	})
	batchErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ingest_batch_errors_total",
		Help: "Total ingest batch handler errors",
	})
)

// ObserveBatch records batch handler duration.
func ObserveBatch(start time.Time) {
	batchLatency.Observe(time.Since(start).Seconds())
}

// IncBatchError increments error counter.
func IncBatchError() {
	batchErrors.Inc()
}

// Handler serves Prometheus metrics on /metrics.
func Handler() http.Handler {
	return promhttp.Handler()
}
