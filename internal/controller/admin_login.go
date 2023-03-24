package controller

import (
	"gateway/internal/dto"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/admin/login", adminLogin.AdminLoginRegister)
}

func (adminlogin *AdminLoginController) AdminLoginRegister(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}
