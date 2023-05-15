package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/metrics"
	"gateway/pkg/log"
	"gateway/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		log.Debug("获取node名称", zap.String("loadBalanceAddr", loadBalanceAddr))

		statusCode := c.GetInt("ErrorCode")
		responseTime := time.Since(startTime).Seconds()
		serviceName := serverInAny.(*enity.ServiceDetail).Info.ServiceName

		if statusCode == response.ServerLimiterAllowErrCode {
			metrics.RecordLimiterMetrics(serviceName, loadBalanceAddr)
		} else if statusCode == response.ClientIPLimiterAllowErrCode {
			metrics.RecordLimiterMetrics(serviceName+"_client", loadBalanceAddr)
		}

		// 更新service负载均衡器请求总数
		metrics.RecordRequestTotalMetrics(serviceName, loadBalanceAddr)
		log.Debug("更新service负载均衡器请求总数", zap.String("serviceName", serviceName), zap.String("loadBalanceAddr", loadBalanceAddr))
		// 更新service响应时间指标
		metrics.RecordResponseTimeMetrics(serviceName, loadBalanceAddr, responseTime)
		log.Debug("更新service响应时间指标", zap.String("serviceName", serviceName), zap.String("loadBalanceAddr", loadBalanceAddr), zap.Float64("responseTime", responseTime))
	}
}
