package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	memoryUsagePercent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_percent",
		Help: "The current memory usage (percentage)",
	})

	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "The current CPU usage (percentage)",
	})

	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of requests",
	}, []string{"name", "node"})

	responseTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "response_time_seconds",
		Help:    "The response time of the application",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	}, []string{"name", "node"})

	limiterCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "limiter_count",
		Help: "The total number of limiter events",
	}, []string{"name", "node"})
)
