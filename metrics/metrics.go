package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of requests",
	}, []string{"name"})

	responseTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "response_time_seconds",
		Help:    "The response time of the application",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	}, []string{"name"})

	errorRate = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "error_rate",
		Help: "The total number of errors occurred",
	}, []string{"name"})

	memoryUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "The current memory usage",
	})

	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "The current CPU usage (percentage)",
	})

	limiterCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "limiter_count",
		Help: "The total number of limiter events",
	}, []string{"name"})

	circuitBreakerCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "circuit_breaker_count",
		Help: "The total number of circuit breaker events",
	}, []string{"name", "state"})
)
