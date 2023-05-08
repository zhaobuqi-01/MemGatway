package middleware

import (
	"fmt"
	"gateway/enity"

	"gateway/pkg/response"
	proxy "gateway/proxy/http_proxy/reverse_proxy"
	"gateway/proxy/pkg"

	"github.com/gin-gonic/gin"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*enity.ServiceDetail)

		// 使用GetLoadBalancer方法获取或创建一个LoadBalance实例
		lb, err := pkg.LoadBalanceTransport.GetLoadBalancer(serviceDetail)
		if err != nil {
			response.ResponseError(c, response.GetLoadBalancerErrCode, err)
			c.Abort()
			return
		}

		// // 使用GetTransportor方法获取或创建一个http.Transport实例
		trans, err := pkg.LoadBalanceTransport.GetTransportor(serviceDetail)
		if err != nil {
			response.ResponseError(c, response.GetTransportorErrCode, err)
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
