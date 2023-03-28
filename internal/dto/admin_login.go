package dto

import (
	"gateway/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	Username string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,is_validate_username)"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456"validate:"required,min=6"`
}

type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` //token
}

func (param *AdminLoginInput) BindValParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}
