package router

import (
	"fmt"
	"net"

	"gateway/proxy/grpc_proxy/reverse_proxy"
	"gateway/proxy/pkg"

	"gateway/proxy/grpc_proxy/middleware"
	"gateway/proxy/grpc_proxy/proxy"

	"gateway/enity"
	"gateway/pkg/log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var grpcServerList = []*warpGrpcServer{}

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcProxyServerRun() {
	serviceList := pkg.Cache.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *enity.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)
			rb, err := pkg.LoadBalanceTransport.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatal("get tcpLoadBalancer failed", zap.String("addr", addr), zap.Error(err))
				return
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatal(" grpcProxy listen failed", zap.String("addr", addr), zap.Error(err))
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
			log.Info("grpcProxy running", zap.String("addr", addr))
			if err := s.Serve(lis); err != nil {
				log.Fatal("grpcProxy fail to run ", zap.String("addr", addr), zap.Error(err))
			}
		}(tempItem)
	}
}

func GrpcProxyServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Info("grpcProxy is stopped", zap.String("addr", grpcServer.Addr))
	}
}
