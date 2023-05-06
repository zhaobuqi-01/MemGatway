package middleware

import (
	"fmt"
	"gateway/dao"
	"gateway/http_proxy/proxy"
	"gateway/utils"

	"github.com/gin-gonic/gin"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			utils.ResponseError(c, utils.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 使用GetLoadBalancer方法获取或创建一个LoadBalance实例
		lb, err := utils.NewLoadBalancer().GetLoadBalancer(serviceDetail)
		if err != nil {
			utils.ResponseError(c, utils.GetLoadBalancerErrCode, err)
			c.Abort()
			return
		}

		// // 使用GetTransportor方法获取或创建一个http.Transport实例
		trans, err := utils.NewTransportor().GetTransportor(serviceDetail)
		if err != nil {
			utils.ResponseError(c, utils.GetTransportorErrCode, err)
			c.Abort()
			return
		}

		//创建 reverseproxy
		//使用 reverseproxy.ServerHTTP(c.Request,c.Response)
		proxy := proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()

		return
	}
}
