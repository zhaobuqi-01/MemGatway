package middleware

import (
	"errors"
	"fmt"
	"gateway/pkg/logger"
	"github.com/spf13/viper"
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
				logger.Panic("panic", zap.String("error", fmt.Sprint(err)), zap.String("stack", string(debug.Stack())))
				if viper.GetString("gin.mode") != "debug" {
					ResponseError(c, 500, errors.New("内部错误"))
					return
				} else {
					ResponseError(c, 500, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}
