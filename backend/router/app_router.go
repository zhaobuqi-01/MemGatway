package router

import (
	"gateway/backend/controller"
	"gateway/backend/middleware"

	"github.com/gin-gonic/gin"
)

func AppRegister(router *gin.Engine) {
	appRouter := router.Group("/app")
	{
		appRouter.Use(
			middleware.SessionAuthMiddleware(),
		)

		controller := controller.NewAPPController()

		appRouter.GET("/app_list", controller.APPList)
		appRouter.GET("/app_detail", controller.APPDetail)
		appRouter.GET("/app_delete", controller.APPDelete)
		appRouter.POST("/app_add", controller.APPAdd)
		appRouter.POST("/app_update", controller.APPUpdate)
		appRouter.GET("/app_stat", controller.APPStat)

	}
}
