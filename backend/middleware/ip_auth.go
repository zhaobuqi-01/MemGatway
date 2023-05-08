package middleware

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/log"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// IPAuthMiddleware IP认证中间件
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取allow_ip配置的string slice类型
		allowIPs := configs.GetGatewayServerConfig().AllowIP
		isMatched := false
		for _, host := range allowIPs {
			// 判断客户端IP是否匹配配置中的IP
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			// 如果不匹配，则返回错误信息
			log.Error("ip not match", zap.String("ip", c.ClientIP()), zap.String("trace_id", c.GetString("TraceID")))
			response.ResponseError(c, response.IpMismatchErrCode, fmt.Errorf("%v, not in iplist", c.ClientIP()))
			c.Abort()
			return
		}
		c.Next()
	}
}
