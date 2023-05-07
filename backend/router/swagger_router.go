package router

import (
	v1 "gateway/backend/api/v1"
	"gateway/configs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func swaggerRegister(router *gin.Engine) {
	// programatically set swagger info
	v1.SwaggerInfo.Title = configs.GetSwaggerConfig().Title
	v1.SwaggerInfo.Description = configs.GetSwaggerConfig().Description
	v1.SwaggerInfo.Version = configs.GetSwaggerConfig().Version
	v1.SwaggerInfo.Host = configs.GetSwaggerConfig().Host
	// v1.SwaggerInfo.BasePath = configs.GetSwaggerConfig().BasePath
	v1.SwaggerInfo.Schemes = configs.GetSwaggerConfig().Schemes
	// Swagger API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
