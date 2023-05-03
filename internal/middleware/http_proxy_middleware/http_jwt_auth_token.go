package http_proxy_middleware

import (
	"fmt"
	"strings"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/pkg"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/gin"
)

// jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			pkg.ResponseError(c, 2001, fmt.Errorf("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		//fmt.Println("token",token)
		appMatched := false
		if token != "" {
			claims, err := public.JwtDecode(token)
			if err != nil {
				pkg.ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			//fmt.Println("claims.Issuer",claims.Issuer)
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			pkg.ResponseError(c, 2003, fmt.Errorf("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}
