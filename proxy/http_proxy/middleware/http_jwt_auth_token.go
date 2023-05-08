package middleware

import (
	"fmt"
	"strings"

	"gateway/enity"
	"gateway/pkg/log"
	"gateway/proxy/http_proxy/utils"
	"gateway/proxy/pkg"

	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*enity.ServiceDetail)

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		log.Debug("token", zap.String("token", token))

		appMatched := false
		if token != "" {
			claims, err := utils.JwtDecode(token)
			if err != nil {
				response.ResponseError(c, response.JwtDecodeErrCode, err)
				c.Abort()
				return
			}
			log.Debug("claims.Issuer", zap.String("claims.Issuer", claims.Issuer))

			appList := pkg.Cache.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			response.ResponseError(c, response.ValidAppErrCode, fmt.Errorf("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}
