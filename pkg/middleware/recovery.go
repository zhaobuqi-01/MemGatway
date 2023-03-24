package middleware

import (
	"fmt"
	"gateway/pkg/logger"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误堆栈和日志
				requestID := c.Writer.Header().Get("X-Request-Id")
				logger.Error(fmt.Sprint(err), zap.String("stack", string(debug.Stack())), zap.String("request_id", requestID))
			}
		}()
		c.Next()
	}
}
