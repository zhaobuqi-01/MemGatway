package controller

import (
	"encoding/json"
	"fmt"
	"gateway/internal/dto"
	"gateway/internal/pkg"
	"gateway/internal/repository"
	"gateway/pkg/logger"
	"gateway/pkg/middleware"
	"time"

	"github.com/gin-gonic/contrib/sessions"
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
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin/login [post]
func (adminlogin *AdminController) AdminLogin(c *gin.Context) {
	var err error
	params := &dto.AdminLoginInput{}
	if err = params.BindValParam(c); err != nil {
		logger.ErrorWithTraceID(c, "parameter binding error")
		middleware.ResponseError(c, 1001, err)
		return
	}

	repo := &repository.AdminRepo{}
	admin, err := repo.LoginCheck(c, params)
	if err != nil {
		logger.ErrorWithTraceID(c, "Login check failed")
		middleware.ResponseError(c, 1002, err)
		return
	}

	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}

	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		logger.ErrorWithTraceID(c, "session serialization failed")
		middleware.ResponseError(c, 1003, err)
		return
	}

	sess := sessions.Default(c)
	sess.Set(pkg.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, "login successful", out)
	logger.InfoWithTraceID(c, "login successful")
}

// AdminLogin godoc
// @Summary 管理员退出登陆
// @Description 管理员退出登陆
// @Tags 管理员接口
// @ID /admin/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/login_out [get]
func (adminlogin *AdminController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(pkg.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "Log out successfully", "")
	logger.InfoWithTraceID(c, "Log out successfully")
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminInfo *AdminController) AdminInfo(c *gin.Context) {
	// 读取seesionKey对应的json字符串转化为结构体
	// 取出数据 封装输出
	sess := sessions.Default(c)
	sessInfo := sess.Get(pkg.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		logger.ErrorWithTraceID(c, "Session deserialization failed")
		middleware.ResponseError(c, 2001, err)
		return
	}

	out := &dto.AminInfoOutput{
		ID:            adminSessionInfo.ID,
		Name:          adminSessionInfo.UserName,
		LoginTime:     adminSessionInfo.LoginTime,
		Avatar:        "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png",
		Introduceions: "I am a super administrator",
		Roles: []string{
			"admin",
		},
	}
	middleware.ResponseSuccess(c, "Obtained administrator information successfully ", out)
	logger.InfoWithTraceID(c, "Obtained administrator information successfully ")
}
