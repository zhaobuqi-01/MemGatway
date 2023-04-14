package router

import (
	"gateway/internal/controller"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ServiceRegister(router *gin.Engine, db *gorm.DB) {
	serviceRouter := router.Group("/service")
	{
		serviceRouter.Use(
			middleware.SessionAuthMiddleware(),
		)

		controller := controller.NewServiceController(db)

		serviceRouter.GET("/service_list", controller.ServiceList)
		serviceRouter.GET("/service_delete", controller.ServiceDelete)
		serviceRouter.POST("/service_add_http", controller.ServiceAddHttp)
		serviceRouter.POST("/service_update_http", controller.ServiceUpdateHttp)
		serviceRouter.GET("/service_detail", controller.ServiceDetail)
	}
}
