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
	"gateway/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 使用默认中间件（logger 和 recovery 中间件）创建 gin 路由
	router := gin.Default()
	router.Use(middleware.SetTraceID)

	// 注册swagger路由
	swaggerRegister(router)

	// 注册admin路由
	AdminRegister(router)

	return router
}
