package http_proxy_middleware

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/pkg"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			pkg.ResponseError(c, pkg.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			log.Info("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := pkg.NewFlowLimiter().GetLimiter(
				pkg.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				pkg.ResponseError(c, pkg.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				pkg.ResponseError(c, pkg.ServerLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit)))
				c.Abort()
				return
			}
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := pkg.NewFlowLimiter().GetLimiter(
				pkg.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				pkg.ResponseError(c, pkg.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				pkg.ResponseError(c, pkg.ClientIPLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("%v flow limit %v", c.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
