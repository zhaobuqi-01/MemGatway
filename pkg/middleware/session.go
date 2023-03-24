package middleware

import (
	"errors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionAuthMiddleware 鉴权中间件，用于检查请求是否携带了正确的 session 信息
func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// 从 session 中获取用户信息，如果用户未登录，则返回错误信息并中止请求
		if name, ok := session.Get("user").(string); !ok || name == "" {
			ResponseError(c, InternalErrorCode, errors.New("user not login")) // 返回未登录错误信息
			c.Abort()                                                         // 中止请求，不再继续执行
			return
		}
		c.Next() // 如果已登录，则继续执行下一个中间件或路由处理函数
	}
}
