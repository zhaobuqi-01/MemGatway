package http_proxy_middleware

import (
	"gateway/internal/dao"
	"gateway/internal/pkg"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
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
