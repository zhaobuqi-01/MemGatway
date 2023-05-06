package router

import (
	"gateway/backend/controller"
	"gateway/backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AppRegister(router *gin.Engine, db *gorm.DB) {
	appRouter := router.Group("/app")
	{
		appRouter.Use(
			middleware.SessionAuthMiddleware(),
		)

		controller := controller.NewAPPController(db)

		appRouter.GET("/app_list", controller.APPList)
		appRouter.GET("/app_detail", controller.APPDetail)
		appRouter.GET("/app_delete", controller.APPDelete)
		appRouter.POST("/app_add", controller.APPAdd)
		appRouter.POST("/app_update", controller.APPUpdate)

	}
}
