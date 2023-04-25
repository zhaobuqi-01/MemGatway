package router

import (
	"gateway/internal/controller"
	"gateway/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRegister(router *gin.Engine, tx *gorm.DB) {

	adminRouter := router.Group("/admin")
	{
		controller := controller.NewAdminController(tx)

		adminRouter.POST("/login", controller.AdminLogin)
		adminRouter.GET("/login_out", controller.AdminLoginOut)
		adminRouter.GET("/admin_info", middleware.SessionAuthMiddleware(), controller.AdminInfo)
		adminRouter.POST("/change_pwd", middleware.SessionAuthMiddleware(), controller.AdminChangePwd)
	}
}
