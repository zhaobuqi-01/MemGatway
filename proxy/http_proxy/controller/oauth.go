package controller

import (
	"gateway/globals"
	"gateway/pkg/log"
	"gateway/pkg/response"
	"gateway/proxy/http_proxy/dto"
	"gateway/proxy/http_proxy/logic"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OAuthController struct {
	aouthLogic logic.OAuthLogic
}

func NewOAuthController() *OAuthController {
	return &OAuthController{aouthLogic: logic.NewOAuthLogic()}
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} response.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oc *OAuthController) Tokens(c *gin.Context) {
	log.Debug("start get Tokens")
	// 参数绑定
	params := &dto.TokensInput{}
	if err := params.BindValParam(c); err != nil {
		return
	}

	out, err := oc.aouthLogic.Tokens(c, params)
	if err != nil {
		response.ResponseError(c, response.TokensErrCode, err)
		log.Error("Failed to get tokens", zap.Error(err))
		return
	}

	response.ResponseSuccess(c, "Get tokens successfully", out)
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID  /admin/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=string} "success"
// @Router /admin/login_out [get]
func (o *OAuthController) AdminLoginOut(c *gin.Context) {
	// 获取session
	session := sessions.Default(c)
	// 删除session
	session.Delete(globals.AdminSessionInfoKey)
	// 保存session
	session.Save()
	response.ResponseSuccess(c, "logout successfully", nil)
}
