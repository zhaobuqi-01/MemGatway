package main

import (
	"context"
	"gateway/configs"
	"gateway/internal/pkg"
	"gateway/internal/router"
	"gateway/pkg/database"
	"gateway/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化资源
	pkg.InitAll()

	// 设置gin模式
	gin.SetMode(configs.GetGinConfig().Mode)

	// 初始化路由
	db := database.GetDB()
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
		logger.Info("server start running", zap.String("addr", serverConfig.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: ", zap.String("addr", serverConfig.Addr), zap.Error(err))
		}
		logger.Info("server is running", zap.String("addr", serverConfig.Addr))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 停止 HTTP 服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:  ", zap.Error(err))
	}
	logger.Info("Server exiting")

	// 等待HttpServerStop执行完毕，执行清理操作
	pkg.CleanupMySQL()
	pkg.CleanupRedis()
	pkg.CleanupLogger()
}
