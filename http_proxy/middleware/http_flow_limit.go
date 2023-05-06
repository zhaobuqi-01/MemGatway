package middleware

import (
	"fmt"

	"gateway/dao"
	"gateway/pkg/log"
	"gateway/utils"

	"github.com/gin-gonic/gin"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			utils.ResponseError(c, utils.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			log.Info("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := utils.NewFlowLimiter().GetLimiter(
				utils.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				utils.ResponseError(c, utils.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				utils.ResponseError(c, utils.ServerLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit)))
				c.Abort()
				return
			}
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := utils.NewFlowLimiter().GetLimiter(
				utils.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				utils.ResponseError(c, utils.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				utils.ResponseError(c, utils.ClientIPLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("%v flow limit %v", c.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
