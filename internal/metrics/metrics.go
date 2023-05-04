package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of requests",
	}, []string{"name", "addr", "type"})

	responseTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "response_time_seconds",
		Help:    "The response time of the application",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	}, []string{"type"})

	throughput = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "throughput",
		Help: "The number of requests processed per second",
	}, []string{"type"})

	errorRate = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "error_rate",
		Help: "The total number of errors occurred",
	}, []string{"type"})

	memoryUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "The current memory usage",
	}, []string{"type"})

	cpuUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "The current CPU usage (percentage)",
	}, []string{"type"})

	limiterCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "limiter_count",
		Help: "The total number of limiter events",
	}, []string{"type"})

	circuitBreakerCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "circuit_breaker_count",
		Help: "The total number of circuit breaker events",
	}, []string{"type", "state"})
)
