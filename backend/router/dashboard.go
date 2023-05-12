package router

import (
	"gateway/backend/controller"
	"gateway/backend/middleware"

	"github.com/gin-gonic/gin"
)

func DashboardRegister(router *gin.Engine) {
	dashboardRouter := router.Group("/dashboard")
	{
		dashboardRouter.Use(
			middleware.SessionAuthMiddleware(),
		)

		controller := controller.NewDashboardController()

		dashboardRouter.GET("/panel_group_data", controller.PanelGroupData)
		dashboardRouter.GET("/service_stat", controller.ServiceStat)
		dashboardRouter.GET("flow_stat", controller.FlowStat)
	}
}
