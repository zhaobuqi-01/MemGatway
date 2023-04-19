package main

import (
	"context"
	"gateway/configs"
	Init "gateway/init"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"gateway/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	Init.InitAll()
	defer Init.CleanupAll()

	// 设置gin模式
	gin.SetMode(configs.GetGinConfig().Mode)

	// 初始化路由
	db := mysql.GetDB()
	r := router.InitRouter(db)

	serverConfig := configs.GetServerConfig()

	srv := &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}
	go func() {
		log.Info("server start running", zap.String("addr", serverConfig.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", zap.String("addr", serverConfig.Addr), zap.Error(err))
		}
		log.Info("server is running", zap.String("addr", serverConfig.Addr))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// 停止 HTTP 服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:  ", zap.Error(err))
	}
	log.Info("Server exiting")
}
