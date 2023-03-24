package dto

import (
	"gateway/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	Username string `json:"username" example:"admin" comment:"用户名" validate:"required"`
	Password string `json:"password" example:"123456" comment:"密码" validate:"required"`
}

func (param *AdminLoginInput) BindValParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}
