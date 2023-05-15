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
		// // 根据 serviceDetail 中的规则判断是否需要协议转换
		// if serviceDetail.HTTPRule != nil && serviceDetail.HTTPRule.NeedHttps == 1 {
		// 	// 进行协议转换，将请求的协议从 http 转换为 https
		// 	if c.Request.URL.Scheme == "http" {
		// 		c.Request.URL.Scheme = "https"
		// 		c.Request.URL.Host = c.Request.Host
		// 	}
		// }

		// // 判断是否需要启用 WebSocket 转换
		// if serviceDetail.HTTPRule != nil && serviceDetail.HTTPRule.NeedWebsocket == 1 {
		// 	// 如果是 WebSocket 请求，则进行 WebSocket 转换
		// 	if isWebSocketRequest(c.Request) {
		// 		handleWebSocketProxy(c)
		// 		return
		// 	}
		// }

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

// func handleWebSocketProxy(c *gin.Context) {
// 	// 使用负载均衡器获取后端服务器实例
// 	serverInterface, ok := c.Get("service")
// 	if !ok {
// 		response.ResponseError(c, response.ServiceNotFoundErrCode, fmt.Errorf("service not found"))
// 		c.Abort()
// 		return
// 	}

// 	serviceDetail := serverInterface.(*enity.ServiceDetail)
// 	lb, err := pkg.LoadBalanceTransport.GetLoadBalancer(serviceDetail)
// 	if err != nil {
// 		response.ResponseError(c, response.GetLoadBalancerErrCode, err)
// 		c.Abort()
// 		return
// 	}

// 	// 创建 WebSocket 转发器
// 	proxy := NewWebSocketReverseProxy(c, lb)
// 	proxy.ServeHTTP(c.Writer, c.Request)
// 	c.Abort()
// }

// func NewWebSocketReverseProxy(c *gin.Context, lb load_balance.LoadBalance) *WebSocketReverseProxy {
// 	// 请求协调者
// 	director := func(req *http.Request) {
// 		nextAddr, err := lb.Get(req.URL.String())
// 		if err != nil || nextAddr == "" {
// 			panic("get next addr fail")
// 		}

// 		target, err := url.Parse(nextAddr)
// 		if err != nil {
// 			log.Error("parse addr fail", zap.Error(err))
// 		}

// 		req.URL.Scheme = target.Scheme
// 		req.URL.Host = target.Host
// 		req.Host = target.Host
// 	}

// 	return &WebSocketReverseProxy{Director: director}
// }

// type WebSocketReverseProxy struct {
// 	Director func(req *http.Request)
// }

// func (p *WebSocketReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	p.Director(r)
// 	proxy := httputil.ReverseProxy{
// 		Director: p.Director,
// 	}

// 	// 检查是否存在 Upgrade 和 Connection 头部
// 	upgradeHeader := r.Header.Get("Upgrade")
// 	connectionHeader := r.Header.Get("Connection")

// 	// 检查是否存在升级标志
// 	isWebSocketRequest := strings.Contains(strings.ToLower(connectionHeader), "upgrade")

// 	// 如果 Upgrade 头部不存在或不为 "websocket"，或者 Connection 头部中没有升级标志，则返回错误
// 	if upgradeHeader != "websocket" || !isWebSocketRequest {
// 		http.Error(w, "WebSocket Upgrade required", http.StatusBadRequest)
// 		return
// 	}

// 	proxy.ServeHTTP(w, r)
// }

// func isWebSocketRequest(r *http.Request) bool {
// 	// 判断请求头中的 Upgrade 字段是否包含 "websocket"
// 	return strings.Contains(strings.ToLower(r.Header.Get("Upgrade")), "websocket")
// }
