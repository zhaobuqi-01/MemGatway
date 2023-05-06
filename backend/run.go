package backend

import (
	"context"
	"gateway/backend/router"
	"gateway/configs"

	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var gatewaySrv *http.Server

func GatewayServerRun() {
	// 初始化路由
	db := mysql.GetDB()
	r := router.InitRouter(db)

	serverConfig := configs.GetGatewayServerConfig()
	log.Info("gatewayServerConfig", zap.Any("serverConfig", serverConfig))

	gatewaySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}

	log.Info("gatewayServer start running", zap.String("addr", serverConfig.Addr))
	if err := gatewaySrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: ", zap.String("gatewayAddr", serverConfig.Addr), zap.Error(err))
	}
	log.Info("gatewayServer is running", zap.String("addr", serverConfig.Addr))

}

func GatewayServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := gatewaySrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:  ", zap.Error(err))
	}
	log.Info("Server exiting")
}
