package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误堆栈和日志
				requestID := c.Writer.Header().Get("X-Request-Id")
				logger.Error(fmt.Sprint(err), zap.String("stack", string(debug.Stack())), zap.String("request_id", requestID))

				// 根据日志级别决定是否将错误信息返回给客户端
				if logger.Core().Enabled(zapcore.ErrorLevel) {
					c.AbortWithStatusJSON(500, gin.H{
						"message": "Internal Server Error",
					})
				} else {
					c.AbortWithStatusJSON(500, gin.H{
						"message": fmt.Sprint(err),
					})
				}
			}
		}()
		c.Next()
	}
}
