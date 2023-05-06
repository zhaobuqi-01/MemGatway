package router

import (
	"gateway/backend/controller"
	"gateway/backend/middleware"

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
		serviceRouter.GET("/service_detail", controller.ServiceDetail)
		serviceRouter.POST("/service_add_http", controller.ServiceAddHttp)
		serviceRouter.POST("/service_update_http", controller.ServiceUpdateHttp)
		serviceRouter.POST("/service_add_tcp", controller.ServiceAddTcp)
		serviceRouter.POST("/service_update_tcp", controller.ServiceUpdateTcp)
		serviceRouter.POST("/service_add_grpc", controller.ServiceAddGrpc)
		serviceRouter.POST("/service_update_grpc", controller.ServiceUpdateGrpc)
	}
}
