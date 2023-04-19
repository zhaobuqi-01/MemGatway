package middleware

import (
	"fmt"
	"gateway/configs"
	"gateway/internal/pkg"
	"gateway/pkg/log"
	"runtime/debug"

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
				log.Error(fmt.Sprint(err), zap.String("stack", string(debug.Stack())), zap.String("request_id", requestID))
				if configs.GetStringConfig("gin.mode") != "debug" {
					pkg.ResponseError(c, pkg.InternalErrorCode, fmt.Errorf("internal error"))
					return
				} else {
					pkg.ResponseError(c, pkg.InternalErrorCode, fmt.Errorf(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}
