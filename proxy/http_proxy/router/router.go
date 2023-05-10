package router

import (
	"context"
	"gateway/configs"
	"gateway/pkg/log"
	"gateway/proxy/http_proxy/controller"
	"gateway/proxy/http_proxy/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(
		middleware.SetTraceID(),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
	)

	// 注册oauth路由
	aouthRouter := router.Group("/oauth")
	aouthRouter.Use(middleware.TranslationMiddleware())
	{
		controller := controller.NewOAuthController()
		aouthRouter.POST("/tokens", controller.Tokens)
	}

	// 注册prometheus监控路由
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 注册健康检查路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Use(
		middleware.HTTPAccessModeMiddleware(),
		middleware.HTTPTrafficStats(),
		middleware.HTTPFlowLimitMiddleware(),
		middleware.HTTPJwtAuthTokenMiddleware(),
		middleware.HTTPJwtFlowCountMiddleware(),
		middleware.HTTPJwtFlowLimitMiddleware(),
		middleware.HTTPWhiteListMiddleware(),
		middleware.HTTPBlackListMiddleware(),
		middleware.HTTPHeaderTransferMiddleware(),
		middleware.HTTPStripUriMiddleware(),
		middleware.HTTPUrlRewriteMiddleware(),
		middleware.HTTPReverseProxyMiddleware(),
	)

	return router
}

var (
	htppsProxySrv *http.Server
	htppProxySrv  *http.Server
)

func HtppProxyServerRun() {
	// 初始化路由
	r := InitRouter()

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
	r := InitRouter()

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
	if err := htppsProxySrv.ListenAndServeTLS("proxy/http_proxy/cert_file/server.crt", "proxy/http_proxy/cert_file/server.key"); err != nil && err != http.ErrServerClosed {
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
