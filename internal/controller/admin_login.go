package controller

import (
	"encoding/json"
	"gateway/internal/common"
	"gateway/internal/dto"
	"gateway/internal/repository"
	"gateway/pkg/logger"
	"gateway/pkg/middleware"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	var err error
	params := &dto.AdminLoginInput{}
	if err = params.BindValParam(c); err != nil {
		logger.ErrorWithTraceID(c, "参数绑定错误")
		middleware.ResponseError(c, 1001, err)
		return
	}

	admin := &repository.Admin{}
	admin, err = admin.LoginCheck(c, params)
	if err != nil {
		logger.ErrorWithTraceID(c, "登录检查失败")
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
		logger.ErrorWithTraceID(c, "session序列化失败")
		middleware.ResponseError(c, 1003, err)
		return
	}

	sess := sessions.Default(c)
	sess.Set(common.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
	logger.InfoWithTraceID(c, "登录成功")
}
