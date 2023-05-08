package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/response"
	"gateway/proxy/pkg"

	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*enity.App)
		if appInfo.Qps > 0 {
			clientLimiter, err := pkg.FlowLimiter.GetLimiter(
				appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				response.ResponseError(c, response.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				response.ResponseError(c, response.APPLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
