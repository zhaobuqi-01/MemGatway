package http_proxy_middleware

import (
	"fmt"
	"strings"

	"gateway/internal/dao"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			pkg.ResponseError(c, pkg.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		iplist := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			if !pkg.InStringSlice(iplist, c.ClientIP()) {
				pkg.ResponseError(c, pkg.ClientIPNotInWhiteListCode, fmt.Errorf(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
