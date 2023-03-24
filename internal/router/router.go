package router

import (
	"gateway/internal/controller"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 设置 Swagger 相关信息
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 创建默认的 gin 实例
	router := gin.Default()

	// 使用传入的中间件
	router.Use(middlewares...)

	// 测试接口
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Swagger 接口
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// demo 接口
	v1 := router.Group("/demo")
	v1.Use(middleware.RecoveryMiddleware(), middleware.RequestLog(), middleware.IPAuthMiddleware(), middleware.TranslationMiddleware())
	{
		controller.DemoRegister(v1)
	}

	// 非登陆接口
	store := sessions.NewCookieStore([]byte("secret"))
	apiNormalGroup := router.Group("/api")
	apiNormalGroup.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.ApiRegister(apiNormalGroup)
	}

	// 登陆接口
	apiAuthGroup := router.Group("/api")
	apiAuthGroup.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.ApiLoginRegister(apiAuthGroup)
	}

	return router
}
