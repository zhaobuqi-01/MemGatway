package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/metrics"
	"gateway/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func HTTPTrafficStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		serverInAny, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		loadBalanceAddr := c.GetString("service_addr")
		statusCode := c.GetInt("ErrorCode")
		responseTime := time.Since(startTime).Seconds()
		serviceName := serverInAny.(*enity.ServiceDetail).Info.ServiceName

		if statusCode == response.ServerLimiterAllowErrCode {
			metrics.RecordLimiterMetrics(serviceName)
		} else if statusCode == response.ClientIPLimiterAllowErrCode {
			metrics.RecordLimiterMetrics(serviceName + "_client")
		}
		// 更新service请求总数
		metrics.RecordRequestTotalMetrics(serviceName)
		// 更新service负载均衡器请求总数
		metrics.RecordRequestTotalMetrics(loadBalanceAddr)
		// 更新service响应时间指标
		metrics.RecordResponseTimeMetrics(serviceName, responseTime)
	}
}

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

		// 如果发生错误，请服务器更新错误率
		if statusCode != response.SuccessCode {
			metrics.RecordErrorRateMetrics("http_proxy")
		}
		// 更新服务器请求总数
		metrics.RecordRequestTotalMetrics("http_proxy")
		// 更新服务器响应时间直方图
		responseTime := time.Since(startTime).Seconds()
		metrics.RecordResponseTimeMetrics("http_proxy", responseTime)
	}
}
