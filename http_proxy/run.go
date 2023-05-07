package http_proxy

import (
	"context"
	"gateway/configs"
	"gateway/http_proxy/router"
	"gateway/pkg/log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	htppsProxySrv *http.Server
	htppProxySrv  *http.Server
)

func HtppProxyServerRun() {
	// 初始化路由
	r := router.InitRouter()

	serverConfig := configs.GetHttpProxyConfig()
	log.Info("httpServerConfig", zap.Any("serverConfig", serverConfig))

	htppProxySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}

	log.Info("HtppProxyServer start running", zap.String("addr", serverConfig.Addr))
	if err := htppProxySrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: ", zap.String("httpProyxAddr", serverConfig.Addr), zap.Error(err))
	}
	log.Info("HtppProxyServer is running", zap.String("addr", serverConfig.Addr))

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
	// 初始化路由
	r := router.InitRouter()

	serverConfig := configs.GetHttpsProxyConfig()
	log.Info("httpsServerConfig", zap.Any("serverConfig", serverConfig))

	htppsProxySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}

	log.Info("HtppsProxyServer start running", zap.String("addr", serverConfig.Addr))
	if err := htppsProxySrv.ListenAndServeTLS("./http_proxy/cert_file/server.crt", "./http_proxy/cert_file/server.key"); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: ", zap.String("httpsProxyAddr", serverConfig.Addr), zap.Error(err))
	}
	log.Info("HtppsProxyServer is running", zap.String("addr", serverConfig.Addr))

}
func HttpsProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := htppsProxySrv.Shutdown(ctx); err != nil {
		log.Fatal("HtppsProxyServer forced to shutdown:  ", zap.Error(err))
	}
	log.Info("HtppsProxyServer exiting")
}
