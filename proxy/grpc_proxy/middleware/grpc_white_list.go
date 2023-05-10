package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/log"
	"gateway/utils"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcWhiteListMiddleware(serviceDetail *enity.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		iplist := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return fmt.Errorf("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			if !utils.InStringSlice(iplist, clientIP) {
				return fmt.Errorf(fmt.Sprintf("%s not in white ip list", clientIP))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Error("RPC failed ", zap.Error(err))
			return err
		}
		return nil
	}
}
