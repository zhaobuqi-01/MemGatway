package middleware

import (
	"gateway/pkg/log"
	"gateway/pkg/response"
	"gateway/proxy/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// serviceManager.HTTPAccessMode 获取service信息
		service, err := pkg.Cache.HTTPAccessMode(c)
		if err != nil {
			// log.Debug("开始记录错误次数")

			// ip := c.ClientIP()
			// // 更新最近的请求时间
			// pkg.RecentRequestTimes.Store(ip, time.Now())
			// // 获取最近的请求时间
			// lastRequestTime, ok := pkg.RecentRequestTimes.Load(ip)
			// if ok && time.Since(lastRequestTime.(time.Time)) < pkg.FrequentRequestTime {
			// 	// 如果最近的请求时间小于设定阈值，增加错误次数
			// 	count, _ := pkg.ErrorCounts.LoadOrStore(ip, 0)
			// 	countVal := count.(int) + 1
			// 	pkg.ErrorCounts.Store(ip, countVal)
			// 	if countVal > pkg.ErrorThreshold {
			// 		pkg.BlackIpCache.Set(ip, true, time.Duration(configs.GetInt("blacklist.expire"))*time.Second)
			// 	}
			// }

			// // 记录错误计数器和blackiplist的数据
			// countVal, _ := pkg.ErrorCounts.Load(ip)
			// log.Error("HTTPAccessModeMiddleware: Error",
			// 	zap.String("IP", ip),
			// 	zap.Int("ErrorCount", countVal.(int)),
			// )

			response.ResponseError(c, response.HTTPAccessModeErrCode, err)
			c.Abort()
			return
		}

		// log记录信息
		log.Info("service info", zap.Any("service info", service))

		c.Set("service", service)
		c.Next()
	}
}
