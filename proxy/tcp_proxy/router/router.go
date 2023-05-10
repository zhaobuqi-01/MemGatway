package router

import (
	"context"
	"fmt"
	"gateway/enity"
	"gateway/proxy/pkg"
	"gateway/proxy/tcp_proxy/middleware"
	"gateway/proxy/tcp_proxy/reverse_proxy"
	"gateway/proxy/tcp_proxy/server"
	"log"
	"net"
)

var tcpServerList = []*server.TcpServer{}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TcpProxyServerRun() {
	serviceList := pkg.Cache.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		if serviceItem == nil {
			continue
		}
		tempItem := serviceItem
		go func(serviceDetail *enity.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)
			rb, err := pkg.LoadBalanceTransport.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			//构建路由及设置中间件
			router := middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				middleware.TCPFlowCountMiddleware(),
				middleware.TCPFlowLimitMiddleware(),
				middleware.TCPWhiteListMiddleware(),
				middleware.TCPBlackListMiddleware(),
			)

			//构建回调handler
			routerHandler := middleware.NewTcpSliceRouterHandler(
				func(c *middleware.TcpSliceRouterContext) server.TCPHandler {
					return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
				}, router)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := &server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != server.ErrServerClosed {
				log.Fatalf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpProxyServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}
