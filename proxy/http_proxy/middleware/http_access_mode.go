package middleware

import (
	"gateway/pkg/log"
	"gateway/pkg/response"
	"gateway/proxy/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// serviceManager.HTTPAccessMode 获取service信息
		service, err := pkg.Cache.HTTPAccessMode(c)
		if err != nil {
			response.ResponseError(c, response.HTTPAccessModeErrCode, err)
			c.Abort()
			return
		}

		// log记录信息
		log.Info("service info", zap.Any("service info", service))

		c.Set("service", service)
		c.Next()
	}
}
