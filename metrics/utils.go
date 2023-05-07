package metrics

import (
	"gateway/pkg/log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

// RecordSystemMetrics records the CPU and memory usage metrics of the system.
func RecordSystemMetrics() {
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
			usedMemoryPercent := (float64(virtMem.Used) / float64(virtMem.Total)) * 100
			memoryUsagePercent.Set(usedMemoryPercent) // 使用 memoryUsagePercent 而不是 memoryUsage
			cpuUsage.Set(cpuPercent[0])               // 不需要乘以 100，因为 gopsutil 已经返回了百分比形式的 CPU 使用率

			time.Sleep(5 * time.Second)
		}
	}()
}

// RecordRequestTotalMetrics records the total number of requests.
func RecordRequestTotalMetrics(serverName string) {
	log.Info("Recording request total metrics", zap.String("server_name", serverName))
	requestsTotal.WithLabelValues(serverName).Inc()
}

// RecordResponseTimeMetrics records the response time of requests.
func RecordResponseTimeMetrics(name string, responseTime float64) {
	responseTimeHistogram.WithLabelValues(name).Observe(responseTime)
}

// RecordErrorRateMetrics records the error rate of requests.
func RecordErrorRateMetrics(name string) {
	errorRate.WithLabelValues(name).Inc()
}

// RecordCircuitBreakerMetrics records the circuit breaker metrics.
func RecordLimiterMetrics(name string) {
	limiterCount.WithLabelValues(name).Inc()
}
