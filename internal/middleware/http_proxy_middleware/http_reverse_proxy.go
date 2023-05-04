package http_proxy_middleware

import (
	"fmt"

	"gateway/internal/dao"
	"gateway/internal/logic"
	"gateway/internal/pkg"
	"gateway/internal/reverse_proxy"

	"github.com/gin-gonic/gin"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			pkg.ResponseError(c, pkg.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 使用GetLoadBalancer方法获取或创建一个LoadBalance实例
		lb, err := logic.NewLoadBalancer().GetLoadBalancer(serviceDetail)
		if err != nil {
			pkg.ResponseError(c, pkg.GetLoadBalancerErrCode, err)
			c.Abort()
			return
		}

		// // 使用GetTransportor方法获取或创建一个http.Transport实例
		trans, err := logic.NewTransportor().GetTransportor(serviceDetail)
		if err != nil {
			pkg.ResponseError(c, pkg.GetTransportorErrCode, err)
			c.Abort()
			return
		}

		//创建 reverseproxy
		//使用 reverseproxy.ServerHTTP(c.Request,c.Response)
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()

		return
	}
}
