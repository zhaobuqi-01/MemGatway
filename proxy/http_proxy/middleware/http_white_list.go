package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/response"
	"gateway/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*enity.ServiceDetail)

		iplist := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			if !utils.InStringSlice(iplist, c.ClientIP()) {
				response.ResponseError(c, response.ClientIPNotInWhiteListCode, fmt.Errorf(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
