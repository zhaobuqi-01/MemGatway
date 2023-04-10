package middleware

import (
	"errors"
	"gateway/internal/pkg"
	"gateway/pkg/logger"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if adminInfo, ok := session.Get(pkg.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
			logger.ErrorWithTraceID(c, "user not login")
			ResponseError(c, InternalErrorCode, errors.New("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
