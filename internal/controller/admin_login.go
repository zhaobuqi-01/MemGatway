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
	group.POST("/admin/login", adminLogin.AdminLogin)
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

// 在AdminLoginController结构体上实现修改密码方法
// AdminLogin godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin_login/change_password
// @Accept  json
// @Produce  json
// @Param body body dto.AdminChangePasswordInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminChangePasswordOutput} "success"
// @Router /admin_login/change_password [post]
// func (adminlogin *AdminLoginController) AdminChangePassword(c *gin.Context) {
// 	params := &dto.AdminChangePasswordInput{}
// 	if err := params.BindValParam(c); err != nil {
// 		middleware.ResponseError(c, 1001, err)
// 		return
// 	}
// 	middleware.ResponseSuccess(c, "")
// }
