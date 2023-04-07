// @title 示例项目 API 文档
// @version 1.0
// @description 这是一个示例项目的 API 文档，包含了项目的所有 API 接口信息。
// @termsOfService https://www.example.com/terms
// @contact.name API 支持团队
// @contact.email support@example.com
// @contact.url https://www.example.com/contact
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
package router

import (
	"fmt"
	v1 "gateway/api/v1"
	"gateway/configs"
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

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 创建默认的 gin 实例
	router := gin.Default()
	router.Use(middleware.SetTraceID)

	// programatically set swagger info
	v1.SwaggerInfo.Title = configs.GetSwaggerConfig().Title
	v1.SwaggerInfo.Description = configs.GetSwaggerConfig().Description
	v1.SwaggerInfo.Version = configs.GetSwaggerConfig().Version
	v1.SwaggerInfo.Host = configs.GetSwaggerConfig().Host
	// v1.SwaggerInfo.BasePath = configs.GetSwaggerConfig().BasePath
	v1.SwaggerInfo.Schemes = configs.GetSwaggerConfig().Schemes
	// Swagger API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// AdminLogin路由
	adminLoginRouter := router.Group("/admin_login")
	fmt.Print(viper.GetString("config.redis.addr"), viper.GetString("config.redis.password"))
	store, err := sessions.NewRedisStore(10, "tcp", configs.GetRedisConfig().Addr, configs.GetRedisConfig().Password, []byte("secret"))
	if err != nil {
		logger.Fatal("sessions.NewRedisSrore err", zap.Error(err))
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

	// AdminInfo路由
	adminRouter := router.Group("/admin")
	adminRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminRegister(adminRouter)
	}

	return router
}
