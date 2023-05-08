package middleware

import (
	"fmt"
	"gateway/pkg/log"
	"gateway/pkg/response"
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
				trace_id := c.GetString("TraceID")
				log.Error(fmt.Sprintf("recover error %v", err), zap.String("trace_id", trace_id))

				response.ResponseError(c, response.InternalErrorCode, fmt.Errorf("internal error"))
				return

			}
		}()
		c.Next()
	}
}
