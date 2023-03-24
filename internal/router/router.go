package router

import (
	"gateway/pkg/logger"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InitRouter 初始化路由，可以传入多个中间件
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// 创建默认的 gin 实例
	router := gin.Default()
	// 使用传入的中间件
	router.Use(middlewares...)

	adminLoginRouter := router.Group("/admin_login")
	store, err := sessions.NewRedisStore(10, "tcp", , []byte("secret"))
	if err !=nil{
		logger.Fatal("sessions.NewRedisSrore err :%v",zap.Error(err))
	}
	adminLoginRouter.Use(
		sessions.Sessions("mysession",store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
		)

	return router
}
