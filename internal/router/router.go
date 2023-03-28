package router

import (
	_ "gateway/api/v1"
	"gateway/internal/controller"
	"gateway/pkg/logger"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Gateway API
// @description API 文档描述
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1
// @Output  /api/v1
// @Summary 初始化路由
// @Description 初始化 gin 实例并注册中间件、API 接口路由等
// @Tags 初始化
// @Accept json
// @Produce json
// @Success 200 {string} string "初始化成功"
// @Router /init [get]
// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 创建默认的 gin 实例
	router := gin.Default()

	// Swagger API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adminLoginRouter := router.Group("/admin_login")
	store, err := sessions.NewRedisStore(10, "tcp", viper.GetString("redis.addr"), "", []byte("secret"))
	if err != nil {
		logger.Fatal("sessions.NewRedisSrore err :%v", zap.Error(err))
	}
	adminLoginRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}

	return router
}
