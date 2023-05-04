package server

import (
	"context"
	"gateway/configs"
	"gateway/internal/metrics"
	"gateway/internal/router/backend_router"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var srv *http.Server

func HtppServerRun() {
	// 设置gin模式
	gin.SetMode(configs.GetGinConfig().Mode)

	// 记录系统指标
	metrics.RecordSystemMetrics("api_gateway")

	// 初始化路由
	db := mysql.GetDB()
	r := backend_router.InitRouter(db)

	serverConfig := configs.GetServerConfig()

	srv = &http.Server{
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
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:  ", zap.Error(err))
	}
	log.Info("Server exiting")
}
