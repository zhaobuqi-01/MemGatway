package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/metrics"
	"gateway/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func HTTPTrafficStats() gin.HandlerFunc {
	return func(c *gin.Context) {

		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*enity.ServiceDetail)

		totalCounter, err := globals.FlowCounter.GetCounter(globals.FlowTotal)
		if err != nil {
			response.ResponseError(c, response.CommErrCode, err)
			c.Abort()
			return
		}
		totalCounter.Increase()
		serviceCounter, err := globals.FlowCounter.GetCounter(serviceDetail.Info.ServiceName)
		if err != nil {
			response.ResponseError(c, response.CommErrCode, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
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
