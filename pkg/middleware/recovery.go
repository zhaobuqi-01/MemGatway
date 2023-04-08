package middleware

import (
	"errors"
	"fmt"
	"gateway/pkg/logger"
	"runtime/debug"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(string(debug.Stack()))
				// 记录错误堆栈和日志
				requestID := c.Writer.Header().Get("X-Request-Id")
				logger.Error(fmt.Sprint(err), zap.String("stack", string(debug.Stack())), zap.String("request_id", requestID))
				if viper.GetString("gin.mode") != "debug" {
					ResponseError(c, InternalErrorCode, errors.New("internal error"))
					return
				} else {
					ResponseError(c, InternalErrorCode, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}
