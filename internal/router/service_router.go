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

		Controller := &controller.ServiceController{}

		serviceRouter.GET("/service_list", Controller.ServiceList)
		serviceRouter.GET("/service_delete", Controller.ServiceDelete)
		serviceRouter.POST("/service_add_http", Controller.ServiceAddHttp)

	}
}
