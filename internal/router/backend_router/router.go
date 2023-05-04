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
package backend_router

import (
	"gateway/configs"
	"gateway/internal/middleware"
	"gateway/pkg/log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()

	store, err := sessions.NewRedisStore(10, "tcp", configs.GetRedisConfig().Addr, configs.GetRedisConfig().Password, []byte("secret"))
	if err != nil {
		log.Fatal("sessions.NewRedisSrore err", zap.Error(err))
	}

	// 注册中间件
	router.Use(
		middleware.SetTraceID(),               // 设置traceID
		sessions.Sessions("mysession", store), // session中间件
		middleware.RecoveryMiddleware(),       // 恢复中间件
		middleware.RequestLog(),               // 请求日志中间件
		middleware.TranslationMiddleware(),    // 国际化中间件
	)

	// 注册prometheus监控路由
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// 注册swagger路由
	swaggerRegister(router)

	// 注册admin路由
	AdminRegister(router, db)

	// 注册service路由
	ServiceRegister(router, db)

	// 注册app路由
	AppRegister(router, db)

	// 注册dashboard路由
	DashboardRegister(router, db)

	return router
}
