package middleware

import (
	"encoding/json"
	"fmt"
	"gateway/enity"
	"gateway/pkg/log"
	"gateway/proxy/pkg"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func GrpcJwtFlowLimitMiddleware(serviceDetail *enity.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return fmt.Errorf("miss metadata from context")
		}
		appInfos := md.Get("app")
		if len(appInfos) == 0 {
			if err := handler(srv, ss); err != nil {
				log.Info("RPC failed with error", zap.Error(err))
				return err
			}
			return nil
		}
		appInfo := &enity.App{}
		if err := json.Unmarshal([]byte(appInfos[0]), appInfo); err != nil {
			return err
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return fmt.Errorf("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		if appInfo.Qps > 0 {
			clientLimiter, err := pkg.FlowLimiter.GetLimiter(
				appInfo.AppID+"_client",
				float64(appInfo.Qps))
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return fmt.Errorf("%v flow limit %v", clientIP, appInfo.Qps)
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Info("RPC failed with error ", zap.Error(err))
			return err
		}
		return nil
	}
}
