package middleware

import (
	"gateway/metrics"
	"gateway/pkg/log"
	"gateway/utils"

	"gateway/dao"

	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appName := appInterface.(*dao.App).Name

		statusCode := c.GetInt("ErrorCode")

		// 记录请求总数
		metrics.RecordRequestTotalMetrics(appName)
		// 记录limit次数
		if statusCode == utils.APPLimiterAllowErrCode {
			metrics.RecordLimiterMetrics(appName)
			log.Info("start app limit")
		}

		c.Next()
	}
}
