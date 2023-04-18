package middleware

import (
	"fmt"
	"gateway/configs"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
)

// IPAuthMiddleware IP认证中间件
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取allow_ip配置的string slice类型
		allowIPs := configs.GetSliceConfig("config.server.allow_ip")
		isMatched := false
		for _, host := range allowIPs {
			// 判断客户端IP是否匹配配置中的IP
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			// 如果不匹配，则返回错误信息
			pkg.ResponseError(c, pkg.IpMismatchErrCode, fmt.Errorf("%v, not in iplist", c.ClientIP()))
			c.Abort()
			return
		}
		c.Next()
	}
}
