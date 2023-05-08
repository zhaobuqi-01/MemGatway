package main

import (
	"context"
	"gateway/backend/router"
	"gateway/backend/utils"
	"gateway/configs"
	"os"
	"os/signal"
	"syscall"

	Init "gateway/init"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 启动后台服务器
	Init.InitAll()
	defer Init.CleanupAll()

	// Create a message queue instance

	db := mysql.GetDB()
	r := router.InitRouter(db)
	utils.InitMq()

	serverConfig := configs.GetGatewayServerConfig()
	log.Info("gatewayServerConfig", zap.Any("serverConfig", serverConfig))

	gatewaySrv := &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}

	go func() {
		log.Info("gatewayServer start running", zap.String("addr", serverConfig.Addr))
		if err := gatewaySrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", zap.String("gatewayAddr", serverConfig.Addr), zap.Error(err))
		}
		log.Info("gatewayServer is running", zap.String("addr", serverConfig.Addr))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := gatewaySrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:  ", zap.Error(err))
	}
	log.Info("Server exiting")

}
