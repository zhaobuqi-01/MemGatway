package middleware

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/log"
	"gateway/proxy/pkg"
	"gateway/utils"

	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// jwt auth token
func GrpcJwtAuthTokenMiddleware(serviceDetail *enity.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return fmt.Errorf("miss metadata from context")
		}
		authToken := ""
		auths := md.Get("authorization")
		if len(auths) > 0 {
			authToken = auths[0]
		}
		token := strings.ReplaceAll(authToken, "Bearer ", "")
		appMatched := false
		if token != "" {
			claims, err := utils.JwtDecode(token)
			if err != nil {
				return fmt.Errorf("JwtDecode %v", err)
			}
			appList := pkg.Cache.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					md.Set("app", utils.Obj2Json(appInfo))
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			return fmt.Errorf("not match valid app")
		}
		if err := handler(srv, ss); err != nil {
			log.Error("grpcJwtAuthTokenMiddleware failed ", zap.Error(err))
			return err
		}
		return nil
	}
}
