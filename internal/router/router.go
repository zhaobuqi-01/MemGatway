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
	v1 "gateway/api/v1"
	"gateway/configs"
	"gateway/internal/controller"
	"gateway/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 使用默认中间件（logger 和 recovery 中间件）创建 gin 路由
	router := gin.Default()
	router.Use(middleware.SetTraceID)

	// 注册swagger路由
	swaggerRegister(router)

	// 注册admin路由
	controller.AdminRegister(router)

	return router
}

func swaggerRegister(router *gin.Engine) {
	// programatically set swagger info
	v1.SwaggerInfo.Title = configs.GetSwaggerConfig().Title
	v1.SwaggerInfo.Description = configs.GetSwaggerConfig().Description
	v1.SwaggerInfo.Version = configs.GetSwaggerConfig().Version
	v1.SwaggerInfo.Host = configs.GetSwaggerConfig().Host
	// v1.SwaggerInfo.BasePath = configs.GetSwaggerConfig().BasePath
	v1.SwaggerInfo.Schemes = configs.GetSwaggerConfig().Schemes
	// Swagger API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
