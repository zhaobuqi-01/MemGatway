package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/log"
	"gateway/pkg/response"
	"gateway/proxy/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			log.Info("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*enity.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := pkg.FlowLimiter.GetLimiter(
				serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				response.ResponseError(c, response.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			log.Info("get serviceLimiter", zap.String("serviceNnme", serviceDetail.Info.ServiceName),
				zap.Any("mysql serviceFlowLimit", serviceDetail.AccessControl.ServiceFlowLimit),
				zap.Any("serviceLimiter", serviceLimiter))
			if !serviceLimiter.Allow() {
				response.ResponseError(c, response.ServerLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit)))
				log.Warn("start server flow limit", zap.Any("service flow limit", serviceDetail.AccessControl.ServiceFlowLimit))
				c.Abort()
				return
			}
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			log.Info("get clientLimiter", zap.Any("serviceNnme", serviceDetail.Info.ServiceName), zap.Any("clientIP", c.ClientIP()))
			clientLimiter, err := pkg.FlowLimiter.GetLimiter(
				serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				response.ResponseError(c, response.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				log.Warn("start client flow limit", zap.String("clientIP", c.ClientIP()), zap.Any("clientIPFlowLimit", serviceDetail.AccessControl.ClientIPFlowLimit))
				response.ResponseError(c, response.ClientIPLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("%v flow limit %v", c.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
