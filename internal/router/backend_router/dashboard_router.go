package backend_router

import (
	"gateway/internal/controller"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DashboardRegister(router *gin.Engine, db *gorm.DB) {
	dashboardRouter := router.Group("/dashboard")
	{
		dashboardRouter.Use(
			middleware.SessionAuthMiddleware(),
		)

		controller := controller.NewDashboardController(db)

		dashboardRouter.GET("/panel_group_data", controller.PanelGroupData)
		dashboardRouter.GET("/service_stat", controller.ServiceStat)
	}
}
