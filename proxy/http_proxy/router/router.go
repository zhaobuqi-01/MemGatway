package router

import (
	"gateway/proxy/http_proxy/controller"
	"gateway/proxy/http_proxy/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(
		middleware.SetTraceID(),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TrafficStats(),
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
		// http_proxy_middleware.HTTPJwtAuthTokenMiddleware(),
		// http_proxy_middleware.HTTPJwtFlowCountMiddleware(),
		// http_proxy_middleware.HTTPJwtFlowLimitMiddleware(),
		middleware.HTTPWhiteListMiddleware(),
		middleware.HTTPBlackListMiddleware(),
		middleware.HTTPHeaderTransferMiddleware(),
		middleware.HTTPStripUriMiddleware(),
		middleware.HTTPUrlRewriteMiddleware(),
		middleware.HTTPReverseProxyMiddleware(),
	)

	return router
}
