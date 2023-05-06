package middleware

import (
	"gateway/metrics"
	"gateway/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func TrafficStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 获取请求信息
		statusCode := c.GetInt("ErrorCode")
		// 更新请求总数
		metrics.RecordRequestTotalMetrics("gateway")

		// 如果发生错误，请更新错误率
		if statusCode != utils.SuccessCode {
			metrics.RecordErrorRateMetrics("gateway")
		}

		// 更新响应时间直方图
		responseTime := time.Since(startTime).Seconds()
		metrics.RecordResponseTimeMetrics("gateway", responseTime)
	}
}
