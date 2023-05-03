package server

import (
	"context"
	"gateway/configs"
	"gateway/internal/router/http_proxy_router"
	"gateway/pkg/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	htppsProxySrv *http.Server
	htppProxySrv  *http.Server
)

func HtppProxyServerRun() {
	// 设置gin模式
	gin.SetMode(configs.GetGinConfig().Mode)

	// 初始化路由
	r := http_proxy_router.InitRouter()

	serverConfig := configs.GetHttpProxyConfig()

	htppProxySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}
	go func() {
		log.Info("HtppProxyServer start running", zap.String("addr", serverConfig.Addr))
		if err := htppProxySrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", zap.String("addr", serverConfig.Addr), zap.Error(err))
		}
		log.Info("HtppProxyServer is running", zap.String("addr", serverConfig.Addr))
	}()
}
func HttpProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := htppProxySrv.Shutdown(ctx); err != nil {
		log.Fatal("HtppProxyServer forced to shutdown:  ", zap.Error(err))
	}
	log.Info("HtppProxyServer exiting")
}

func HttpsProxyServerRun() {
	// 设置gin模式
	gin.SetMode(configs.GetGinConfig().Mode)

	// 初始化路由
	r := http_proxy_router.InitRouter()

	serverConfig := configs.GetHttpsProxyConfig()

	htppsProxySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}
	go func() {
		log.Info("HtppsProxyServer start running", zap.String("addr", serverConfig.Addr))
		if err := htppsProxySrv.ListenAndServeTLS("./cert_file/server.crt", "./cert_file/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", zap.String("addr", serverConfig.Addr), zap.Error(err))
		}
		log.Info("HtppsProxyServer is running", zap.String("addr", serverConfig.Addr))
	}()
}
func HttpsProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := htppsProxySrv.Shutdown(ctx); err != nil {
		log.Fatal("HtppsProxyServer forced to shutdown:  ", zap.Error(err))
	}
	log.Info("HtppsProxyServer exiting")
}
