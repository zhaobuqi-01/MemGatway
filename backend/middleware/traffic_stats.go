package middleware

import (
	"gateway/metrics"
	"gateway/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func TrafficStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果是 Prometheus 请求，不记录指标
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		startTime := time.Now()

		// 处理请求
		c.Next()

		// 获取请求信息
		statusCode := c.GetInt("ErrorCode")

		// 如果发生错误，请更新错误率
		if statusCode != response.SuccessCode {
			metrics.RecordErrorRateMetrics("gateway")
		}
		// 更新请求总数
		metrics.RecordRequestTotalMetrics("gateway")
		// 更新响应时间直方图
		responseTime := time.Since(startTime).Seconds()
		metrics.RecordResponseTimeMetrics("gateway", responseTime)
	}
}
