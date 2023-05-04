package http_proxy_middleware

import (
	"gateway/internal/logic"
	"gateway/internal/pkg"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// serviceManager.HTTPAccessMode 获取service信息
		service, err := logic.NewServiceManager().HTTPAccessMode(c)
		if err != nil {
			pkg.ResponseError(c, pkg.HTTPAccessModeErrCode, err)
			c.Abort()
			return
		}

		// log记录信息
		log.Info("service info", zap.Any("service info", service))

		c.Set("service", service)
		c.Next()
	}
}
