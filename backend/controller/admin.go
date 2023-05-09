package controller

import (
	"gateway/backend/dto"
	"gateway/backend/logic"
	"gateway/pkg/log"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Admin interface {
	AdminLogin(c *gin.Context)
	AdminLoginOut(c *gin.Context)
	AdminInfo(c *gin.Context)
	AdminChangePwd(c *gin.Context)
}

type adminController struct {
	logic.AdminLogic
}

func NewAdminController(db *gorm.DB) *adminController {
	return &adminController{
		logic.NewAdminLogic(db),
	}
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags Admin
// @ID /admin/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} response.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin/login [post]
func (a *adminController) AdminLogin(c *gin.Context) {
	// 参数绑定
	params := &dto.AdminLoginInput{}
	if err := params.BindValParam(c); err != nil {
		return
	}

	sessInfo, err := a.Login(c, params)
	if err != nil {
		response.ResponseError(c, response.UserLoginErrCode, err)
		log.Error("Admin login failed", zap.String("username", params.UserName), zap.Error(err))
		return

	}

	out := &dto.AdminLoginOutput{Token: sessInfo.UserName}
	response.ResponseSuccess(c, "login successful", out)
}

// AdminLogin godoc
// @Summary 管理员退出登陆
// @Description 管理员退出登陆
// @Tags Admin
// @ID /admin/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=string} "success"
// @Router /admin/login_out [get]
func (a *adminController) AdminLoginOut(c *gin.Context) {
	err := a.AdminLogout(c)
	if err != nil {
		response.ResponseError(c, response.UserLoginOutErrCode, err)
		log.Error("Admin login out failed", zap.Error(err))
		return
	}

	response.ResponseSuccess(c, "Log out successfully", "")
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags Admin
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=dto.AminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (a *adminController) AdminInfo(c *gin.Context) {
	out, err := a.GetAdminInfo(c)
	if err != nil {
		response.ResponseError(c, response.UserInfoErrCode, err)
		log.Error("Get admin info failed", zap.Error(err))
		return
	}

	response.ResponseSuccess(c, "Obtained administrator information successfully ", out)
}

// AdminInfo godoc
// @Summary 管理员密码修改
// @Description 管理员密码修改
// @Tags Admin
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.AdminChangePwdInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *adminController) AdminChangePwd(c *gin.Context) {
	params := &dto.AdminChangePwdInput{}
	if err := params.BindValParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := a.ChangeAdminPassword(c, params)
	if err != nil {
		response.ResponseError(c, response.UserChangePwdErrCode, err)
		log.Error("Change admin password failed", zap.Error(err))
		return
	}

	response.ResponseSuccess(c, "Password changed successfully", "")
}
