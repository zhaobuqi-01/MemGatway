package router

import (
	"gateway/internal/controller"
	"gateway/internal/middleware"
	"gateway/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRegister(router *gin.Engine, db *gorm.DB) {

	adminRouter := router.Group("/admin")
	{

		db := database.GetDB()
		controller := controller.NewAdminController(db)

		adminRouter.POST("/login", controller.AdminLogin)
		adminRouter.GET("/login_out", controller.AdminLoginOut)
		adminRouter.GET("/admin_info", middleware.SessionAuthMiddleware(), controller.AdminInfo)
		adminRouter.POST("/change_pwd", middleware.SessionAuthMiddleware(), controller.AdminChangePwd)
	}
}
