package middleware

import (
	"gateway/pkg/log"
	"gateway/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// serviceManager.HTTPAccessMode 获取service信息
		service, err := utils.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			utils.ResponseError(c, utils.HTTPAccessModeErrCode, err)
			c.Abort()
			return
		}

		// log记录信息
		log.Info("service info", zap.Any("service info", service))

		c.Set("service", service)
		c.Next()
	}
}
