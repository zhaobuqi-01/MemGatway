package middleware

import (
	"fmt"

	"strings"

	"gateway/dao"
	"gateway/pkg/log"
	"gateway/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("start sBlackListMiddleware")
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
			log.Debug("get whileIpList", zap.Any("whileIpList", whileIpList))
		}

		blackIpList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
			log.Debug("get blackIpList", zap.Any("blackIpList", blackIpList))
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if utils.InStringSlice(blackIpList, c.ClientIP()) {
				utils.ResponseError(c, utils.ClientIPInBlackListErrCode, fmt.Errorf(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				log.Info("client ip in black list", zap.String("clientIP", c.ClientIP()))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
