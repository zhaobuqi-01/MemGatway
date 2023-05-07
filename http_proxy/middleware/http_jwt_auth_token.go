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

// jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			utils.ResponseError(c, utils.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		log.Debug("token", zap.String("token", token))

		appMatched := false
		if token != "" {
			claims, err := utils.JwtDecode(token)
			if err != nil {
				utils.ResponseError(c, utils.JwtDecodeErrCode, err)
				c.Abort()
				return
			}
			log.Debug("claims.Issuer", zap.String("claims.Issuer", claims.Issuer))

			appList := utils.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			utils.ResponseError(c, utils.ValidAppErrCode, fmt.Errorf("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}
