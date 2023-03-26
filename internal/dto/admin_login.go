package dto

import (
	"gateway/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	Username string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456"validate:"required,min=6"`
}

func (param *AdminLoginInput) BindValParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}
