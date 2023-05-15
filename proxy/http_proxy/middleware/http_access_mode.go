package middleware

import (
	"fmt"
	"gateway/pkg/log"
	"gateway/pkg/response"
	"gateway/proxy/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTPAccessModeMiddleware http access mode middleware
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := pkg.Cache.HTTPAccessMode(c)
		if err != nil {
			response.ResponseError(c, response.HTTPAccessModeErrCode, err)
			c.Abort()
			return
		}

		if service == nil {
			response.ResponseError(c, response.HTTPAccessModeErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		log.Info("service info", zap.Any("service info", service))

		c.Set("service", service)
		c.Next()
	}
}
