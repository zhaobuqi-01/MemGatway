package grpc_proxy_router

import (
	"fmt"
	"log"
	"net"

	"gateway/proxy/grpc_proxy/reverse_proxy"
	"gateway/proxy/pkg"

	"gateway/proxy/grpc_proxy/middleware"
	"gateway/proxy/grpc_proxy/proxy"

	"gateway/enity"

	"google.golang.org/grpc"
)

var grpcServerList = []*warpGrpcServer{}

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServerRun() {
	serviceList := pkg.Cache.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *enity.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)
			rb, err := pkg.LoadBalanceTransport.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [INFO] GrpcListen %v err:%v\n", addr, err)
			}
			grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					// middleware.GrpcFlowCountMiddleware(serviceDetail),
					middleware.GrpcFlowLimitMiddleware(serviceDetail),
					middleware.GrpcJwtAuthTokenMiddleware(serviceDetail),
					// middleware.GrpcJwtFlowCountMiddleware(serviceDetail),
					middleware.GrpcJwtFlowLimitMiddleware(serviceDetail),
					middleware.GrpcWhiteListMiddleware(serviceDetail),
					middleware.GrpcBlackListMiddleware(serviceDetail),
					middleware.GrpcHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler))

			grpcServerList = append(grpcServerList, &warpGrpcServer{
				Addr:   addr,
				Server: s,
			})
			log.Printf(" [INFO] grpc_proxy_run %v\n", addr)
			if err := s.Serve(lis); err != nil {
				log.Fatalf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
