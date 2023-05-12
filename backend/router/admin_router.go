package router

import (
	"gateway/backend/controller"
	"gateway/backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRegister(router *gin.Engine) {

	adminRouter := router.Group("/admin")
	{
		c := controller.NewAdminController()

		adminRouter.POST("/login", c.AdminLogin)
		adminRouter.GET("/login_out", c.AdminLoginOut)
		adminRouter.GET("/admin_info", middleware.SessionAuthMiddleware(), c.AdminInfo)
		adminRouter.POST("/change_pwd", middleware.SessionAuthMiddleware(), c.AdminChangePwd)
	}
}
