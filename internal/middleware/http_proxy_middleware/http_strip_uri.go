package http_proxy_middleware

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			pkg.ResponseError(c, pkg.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		if serviceDetail.HTTPRule.RuleType == pkg.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri == 1 {

			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)

		}

		c.Next()
	}
}
