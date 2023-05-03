package http_proxy_middleware

import (
	"fmt"

	"gateway/internal/logic"
	"gateway/internal/reverse_proxy"

	"gateway/internal/dao"

	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			pkg.ResponseError(c, pkg.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 获取LoadBalancer实例
		loadBalancer := logic.NewLoadBalancer()
		// 使用GetLoadBalancer方法获取或创建一个LoadBalance实例
		lb, err := loadBalancer.GetLoadBalancer(serviceDetail)
		if err != nil {
			pkg.ResponseError(c, pkg.GetLoadBalancerErrCode, err)
			c.Abort()
			return
		}

		// 获取Transportor实例
		transportor := logic.NewTransportor()
		// // 使用GetTransportor方法获取或创建一个http.Transport实例
		trans, err := transportor.GetTransportor(serviceDetail)
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
