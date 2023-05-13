package middleware

import (
	"fmt"
	"gateway/configs"
	"gateway/enity"
	"gateway/pkg/log"
	"gateway/proxy/pkg"
	"gateway/utils"
	"strings"
	"time"

	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("start BlackListMiddleware")
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*enity.ServiceDetail)

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
				response.ResponseError(c, response.ClientIPInBlackListErrCode, fmt.Errorf(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				log.Info("client ip in black list", zap.String("clientIP", c.ClientIP()))
				c.Abort()
				return
			}
		}

		if _, found := pkg.BlackIpCache.Get(c.ClientIP()); found {
			response.ResponseError(c, response.ClientIPInBlackListErrCode, fmt.Errorf(fmt.Sprintf("%s in black ip list", c.ClientIP())))
			c.Abort()
			return
		}

		c.Next()

		switch c.GetInt("ErrorCode") {
		case response.AppNotFoundErrCode, response.ServiceNotFoundErrCode, response.HTTPAccessModeErrCode:
			ip := c.ClientIP()
			// 更新最近的请求时间
			pkg.RecentRequestTimes.Store(ip, time.Now())
			// 获取最近的请求时间
			lastRequestTime, ok := pkg.RecentRequestTimes.Load(ip)
			if ok && time.Since(lastRequestTime.(time.Time)) < pkg.FrequentRequestTime {
				// 如果最近的请求时间小于设定阈值，增加错误次数
				count, _ := pkg.ErrorCounts.LoadOrStore(ip, 0)
				count = count.(int) + 1
				pkg.ErrorCounts.Store(ip, count)
				if count.(int) > pkg.ErrorThreshold {
					pkg.BlackIpCache.Set(ip, true, time.Duration(configs.GetInt("blacklist.expire"))*time.Second)
				}
			}
		}
	}
}
