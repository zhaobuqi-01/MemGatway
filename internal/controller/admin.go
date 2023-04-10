package controller

import (
	"gateway/internal/dto"
	"gateway/internal/pkg"
	"gateway/internal/service"
	"gateway/pkg/database"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} pkg.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin/login [post]
func (adminlogin *AdminController) AdminLogin(c *gin.Context) {
	// 参数绑定
	params := &dto.AdminLoginInput{}
	if err := params.BindValParam(c); err != nil {
		logger.ErrorWithTraceID(c, "parameter binding error")
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}

	db := database.GetDB()
	adminService := service.NewAdminService(db)

	sessInfo, err := adminService.Login(c, params)
	if err != nil {
		logger.ErrorWithTraceID(c, "Login check failed")
		pkg.ResponseError(c, pkg.BusinessLogicError, err)
		return
	}

	out := &dto.AdminLoginOutput{Token: sessInfo.UserName}
	pkg.ResponseSuccess(c, "login successful", out)
	logger.InfoWithTraceID(c, "login successful")
}

// AdminLogin godoc
// @Summary 管理员退出登陆
// @Description 管理员退出登陆
// @Tags 管理员接口
// @ID /admin/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /admin/login_out [get]
func (adminloginout *AdminController) AdminLoginOut(c *gin.Context) {

	adminService := service.NewAdminService(nil)
	err := adminService.AdminLogout(c)
	if err != nil {
		pkg.ResponseError(c, pkg.BusinessLogicError, err)
		return
	}

	pkg.ResponseSuccess(c, "Log out successfully", "")
	logger.InfoWithTraceID(c, "Log out successfully")
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response{data=dto.AminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminInfo *AdminController) AdminInfo(c *gin.Context) {
	adminService := service.NewAdminService(nil)
	out, err := adminService.GetAdminInfo(c)
	if err != nil {
		pkg.ResponseError(c, pkg.BusinessLogicError, err)
		return
	}

	pkg.ResponseSuccess(c, "Obtained administrator information successfully ", out)
	logger.InfoWithTraceID(c, "Obtained administrator information successfully ")
}

// AdminInfo godoc
// @Summary 管理员密码修改
// @Description 管理员密码修改
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.AdminChangePwdInput true "body"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminChangePwd *AdminController) AdminChangePwd(c *gin.Context) {
	params := &dto.AdminChangePwdInput{}
	if err := params.BindValParam(c); err != nil {
		logger.ErrorWithTraceID(c, "parameter binding error")
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}

	db := database.GetDB()
	adminService := service.NewAdminService(db)

	err := adminService.ChangeAdminPassword(c, params)
	if err != nil {
		logger.ErrorWithTraceID(c, "Password modification failed")
		pkg.ResponseError(c, pkg.BusinessLogicError, err)
		return
	}

	pkg.ResponseSuccess(c, "Password modification successful", "")
	logger.InfoWithTraceID(c, "Password modification successful")
}
