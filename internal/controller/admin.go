package controller

import (
	"gateway/internal/dto"
	"gateway/internal/logic"
	"gateway/internal/pkg"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminController struct {
	logic logic.AdminLogic
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{
		logic: logic.NewAdminLogic(db),
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
// @Success 200 {object} pkg.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin/login [post]
func (a *AdminController) AdminLogin(c *gin.Context) {
	// 参数绑定
	params := &dto.AdminLoginInput{}
	if err := params.BindValParam(c); err != nil {
		return
	}

	sessInfo, err := a.logic.Login(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		log.Error("Admin login failed", zap.String("username", params.UserName), zap.Error(err))
		return

	}

	out := &dto.AdminLoginOutput{Token: sessInfo.UserName}
	pkg.ResponseSuccess(c, "login successful", out)
}

// AdminLogin godoc
// @Summary 管理员退出登陆
// @Description 管理员退出登陆
// @Tags Admin
// @ID /admin/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /admin/login_out [get]
func (a *AdminController) AdminLoginOut(c *gin.Context) {
	err := a.logic.AdminLogout(c)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		log.Error("Admin login out failed", zap.Error(err))
		return
	}

	pkg.ResponseSuccess(c, "Log out successfully", "")
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags Admin
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response{data=dto.AminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (a *AdminController) AdminInfo(c *gin.Context) {
	out, err := a.logic.GetAdminInfo(c)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		log.Error("Get admin info failed", zap.Error(err))
		return
	}

	pkg.ResponseSuccess(c, "Obtained administrator information successfully ", out)
}

// AdminInfo godoc
// @Summary 管理员密码修改
// @Description 管理员密码修改
// @Tags Admin
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.AdminChangePwdInput true "body"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *AdminController) AdminChangePwd(c *gin.Context) {
	params := &dto.AdminChangePwdInput{}
	if err := params.BindValParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}

	err := a.logic.ChangeAdminPassword(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		log.Error("Change admin password failed", zap.Error(err))
		return
	}

	pkg.ResponseSuccess(c, "Password modification successful", "")
}
