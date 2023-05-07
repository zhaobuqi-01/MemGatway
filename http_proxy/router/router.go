package router

import (
	"gateway/http_proxy/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter(serverName string) *gin.Engine {
	router := gin.Default()

	router.Use(
		middleware.SetTraceID(),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
	)

	// 注册prometheus监控路由
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// 注册健康检查路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// oauth := router.Group("/oauth")
	// oauth.Use(middleware.TranslationMiddleware())

	// 	controller.OAuthRegister(oauth)

	router.Use(
		middleware.HTTPAccessModeMiddleware(),
		middleware.HTTPTrafficStats(serverName),
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
