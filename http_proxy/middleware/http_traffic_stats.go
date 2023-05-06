package middleware

import (
	"gateway/dao"
	"gateway/metrics"
	"gateway/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const proxyServer = "https_proxy/http_proxy"

func HTTPTrafficStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		service, exists := c.Get("service")
		loadBalanceAddr := c.GetString("serverAddr")
		statusCode := c.GetInt("ErrorCode")
		responseTime := time.Since(startTime).Seconds()

		if exists {
			serviceName := service.(*dao.ServiceDetail).Info.ServiceName

			if statusCode == utils.SuccessCode {
				// 更新http/https代理服务器请求总数
				metrics.RecordRequestTotalMetrics(proxyServer)

				// 更新http/https代理服务器响应时间指标
				metrics.RecordResponseTimeMetrics(proxyServer, responseTime)

				// 更新服务指标
				metrics.RecordRequestTotalMetrics(serviceName)
				// 更新服务负载均衡器请求总数
				metrics.RecordRequestTotalMetrics(loadBalanceAddr)
				// 更新服务响应时间指标
				metrics.RecordResponseTimeMetrics(serviceName, responseTime)
			} else {

				// 更新http/https代理服务器错误率
				metrics.RecordErrorRateMetrics(proxyServer)

				// 更新服务错误率
				metrics.RecordErrorRateMetrics(serviceName)
				switch statusCode {
				case utils.ServerLimiterAllowErrCode:
					// 更新服务限流器错误率
					metrics.RecordLimiterMetrics(serviceName)
				case utils.ClientIPLimiterAllowErrCode:
					// 更新服务客户端IP限流器错误率
					metrics.RecordLimiterMetrics(serviceName + "_" + c.ClientIP())
				}
			}

		}
	}
}
