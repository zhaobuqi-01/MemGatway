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

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 创建默认的 gin 实例
	router := gin.Default()

	// programatically set swagger info
	v1.SwaggerInfo.Title = configs.GetSwaggerConfig().Title
	v1.SwaggerInfo.Description = configs.GetSwaggerConfig().Description
	v1.SwaggerInfo.Version = configs.GetSwaggerConfig().Version
	v1.SwaggerInfo.Host = configs.GetSwaggerConfig().Host
	v1.SwaggerInfo.BasePath = configs.GetSwaggerConfig().BasePath
	v1.SwaggerInfo.Schemes = configs.GetSwaggerConfig().Schemes
	// Swagger API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	return router
}
