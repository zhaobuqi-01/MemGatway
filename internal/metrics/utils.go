package metrics

import (
	"gateway/pkg/log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

// RecordSystemMetrics records the CPU and memory usage metrics of the system.
func RecordSystemMetrics(types string) {
	go func() {
		for {
			// Use gopsutil to get accurate CPU usage
			cpuPercent, err := cpu.Percent(time.Second, false)
			if err != nil {
				log.Warn("get cpu percent failed", zap.Error(err))
				continue
			}

			// Use gopsutil to get accurate memory usage
			virtMem, err := mem.VirtualMemory()
			if err != nil {
				log.Warn("get virtual memory failed", zap.Error(err))
				continue
			}
			memoryUsage.WithLabelValues(types).Set(float64(virtMem.Used))
			cpuUsage.WithLabelValues(types).Set(cpuPercent[0])

			time.Sleep(5 * time.Second)
		}
	}()
}

// RecordRequestTotalMetrics records the total number of requests.
func RecordRequestTotalMetrics(serverName, serverAddr string) {
	requestsTotal.WithLabelValues(serverName, serverAddr, "second").Inc()
	requestsTotal.WithLabelValues(serverName, serverAddr, "hour").Inc()
	requestsTotal.WithLabelValues(serverName, serverAddr, "day").Inc()
}

// RecordThroughputMetrics records the number of requests processed per second.
func RecordThroughputMetrics(types string, responseTime time.Duration) {
	throughput.WithLabelValues(types).Set(1 / responseTime.Seconds())
}

// RecordResponseTimeMetrics records the response time of requests.
func RecordResponseTimeMetrics(types string, responseTime time.Duration) {
	responseTimeHistogram.WithLabelValues("api_gateway").Observe(responseTime.Seconds())
}

// RecordErrorRateMetrics records the error rate of requests.
func RecordErrorRateMetrics(types string) {
	errorRate.WithLabelValues(types).Inc()
}

// RecordCircuitBreakerMetrics records the circuit breaker metrics.
func RecordLimiterMetrics(types string) {
	limiterCount.WithLabelValues(types).Inc()
}
