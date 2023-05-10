package dto

import (
	"gateway/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   //密码
}

func (param *AdminLoginInput) BindValParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}

type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` //token
}

type AminInfoOutput struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	LoginTime     time.Time `json:"login_time"`
	Avatar        string    `json:"avatar"`
	Introduceions string    `json:"introduceion"`
	Roles         []string  `json:"roles"`
}

type AdminSessionInfo struct {
	ID        int       `json:"id" form:"id" comment:"管理员ID" example:"1" validate:""`                                         //管理员ID
	UserName  string    `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //管理员用户名
	LoginTime time.Time `json:"login_time" form:"login_time" comment:"登录时间" example:"2020-10-10 10:10:10" validate:""`        //登录时间
}
type AdminChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (param *AdminChangePwdInput) BindValParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}
