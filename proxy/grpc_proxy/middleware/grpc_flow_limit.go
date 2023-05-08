package middleware

import (
	"fmt"
	"gateway/proxy/pkg"
	"log"
	"strings"

	"gateway/enity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcFlowLimitMiddleware(serviceDetail *enity.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := pkg.FlowLimiter.GetLimiter(
				serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				return err
			}
			if !serviceLimiter.Allow() {
				return fmt.Errorf(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit))
			}
		}
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return fmt.Errorf("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := pkg.FlowLimiter.GetLimiter(
				serviceDetail.Info.ServiceName+"_client",
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return fmt.Errorf(fmt.Sprintf("%v flow limit %v", clientIP, serviceDetail.AccessControl.ClientIPFlowLimit))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
