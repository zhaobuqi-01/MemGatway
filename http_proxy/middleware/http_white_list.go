package middleware

import (
	"fmt"
	"gateway/dao"
	"gateway/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			utils.ResponseError(c, utils.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		iplist := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			if !utils.InStringSlice(iplist, c.ClientIP()) {
				utils.ResponseError(c, utils.ClientIPNotInWhiteListCode, fmt.Errorf(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
