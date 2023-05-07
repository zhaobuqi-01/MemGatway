package middleware

import (
	"fmt"
	"gateway/dao"
	"gateway/utils"

	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		if appInfo.Qps > 0 {
			clientLimiter, err := utils.GloablFlowLimiter.GetLimiter(
				utils.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				utils.ResponseError(c, utils.GetLimiterErrCode, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				utils.ResponseError(c, utils.APPLimiterAllowErrCode, fmt.Errorf(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
