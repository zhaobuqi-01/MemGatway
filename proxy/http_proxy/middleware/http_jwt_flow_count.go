package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*enity.App)

		appCounter, err := globals.FlowCounter.GetCounter(appInfo.AppID)
		if err != nil {
			response.ResponseError(c, response.CommErrCode, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.QPD > appInfo.Qpd {
			response.ResponseError(c, 2003, fmt.Errorf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.QPD))
			c.Abort()
			return
		}
		c.Next()

		c.Next()
	}
}
