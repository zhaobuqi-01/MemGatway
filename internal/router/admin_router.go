package router

import (
	"gateway/configs"
	"gateway/internal/controller"
	"gateway/pkg/logger"
	"gateway/pkg/middleware"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdminRegister(router *gin.Engine) {
	store, err := sessions.NewRedisStore(10, "tcp", configs.GetRedisConfig().Addr, configs.GetRedisConfig().Password, []byte("secret"))
	if err != nil {
		logger.Fatal("sessions.NewRedisSrore err", zap.Error(err))
	}

	adminRouter := router.Group("/admin")
	{
		adminRouter.Use(
			sessions.Sessions("mysession", store),
			middleware.RecoveryMiddleware(),
			middleware.RequestLog(),
			middleware.TranslationMiddleware(),
		)

		Controller := &controller.AdminController{}

		adminRouter.POST("/login", Controller.AdminLogin)
		adminRouter.GET("/login_out", Controller.AdminLoginOut)
		adminRouter.GET("/admin_info", middleware.SessionAuthMiddleware(), Controller.AdminInfo)
		adminRouter.POST("/change_pwd", middleware.SessionAuthMiddleware(), Controller.AdminChangePwd)

	}
}
