package router

import (
	"gateway/internal/controller"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRegister(router *gin.Engine) {

	adminRouter := router.Group("/admin")
	{

		Controller := &controller.AdminController{}

		adminRouter.POST("/login", Controller.AdminLogin)
		adminRouter.GET("/login_out", Controller.AdminLoginOut)
		adminRouter.GET("/admin_info", middleware.SessionAuthMiddleware(), Controller.AdminInfo)
		adminRouter.POST("/change_pwd", middleware.SessionAuthMiddleware(), Controller.AdminChangePwd)

	}
}
