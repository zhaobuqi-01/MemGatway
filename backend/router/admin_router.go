package router

import (
	"gateway/backend/controller"
	"gateway/backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRegister(router *gin.Engine, tx *gorm.DB) {

	adminRouter := router.Group("/admin")
	{
		c := controller.NewAdminController(tx)

		adminRouter.POST("/login", c.AdminLogin)
		adminRouter.GET("/login_out", c.AdminLoginOut)
		adminRouter.GET("/admin_info", middleware.SessionAuthMiddleware(), c.AdminInfo)
		adminRouter.POST("/change_pwd", middleware.SessionAuthMiddleware(), c.AdminChangePwd)
	}
}
