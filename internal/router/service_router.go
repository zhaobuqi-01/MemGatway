package router

import (
	"gateway/internal/controller"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServiceRegister(router *gin.Engine) {
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
