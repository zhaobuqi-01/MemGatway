package middleware

import (
	"fmt"

	"strings"

	"gateway/dao"
	"gateway/utils"

	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			utils.ResponseError(c, utils.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		whileIpList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		blackIpList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if utils.InStringSlice(blackIpList, c.ClientIP()) {
				utils.ResponseError(c, utils.ClientIPInBlackListErrCode, fmt.Errorf(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
