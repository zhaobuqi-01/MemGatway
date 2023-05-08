package middleware

import (
	"fmt"

	"gateway/globals"
	"gateway/pkg/log"
	"gateway/pkg/response"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if adminInfo, ok := session.Get(globals.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
			log.Error("user not login", zap.String("trace_id", c.GetString("TraceID")))
			response.ResponseError(c, response.UserNotLoggedInErrCode, fmt.Errorf("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
