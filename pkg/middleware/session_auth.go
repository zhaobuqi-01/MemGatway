package middleware

import (
	"fmt"
	"gateway/internal/pkg"
	"gateway/pkg/log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if adminInfo, ok := session.Get(pkg.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
			log.Error("user not login", zap.String("trace_id", c.GetString("TraceID")))
			pkg.ResponseError(c, pkg.UserNotLoggedInErrCode, fmt.Errorf("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
