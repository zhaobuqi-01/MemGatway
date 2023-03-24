package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type AdminLoginInput struct {
	Username string `json:"username" example:"admin" comment:"用户名" validate:"required"`
	Password string `json:"password" example:"123456" comment:"密码" validate:"required"`
}

func (param *AdminLoginInput) BindValParam(c *gin.Context) error {
	// 创建指针变量 params，并进行绑定和校验
	params := &AdminLoginInput{}
	if err := c.Bind(params); err != nil {
		return err
	}
	if err := c.ShouldBind(params); err != nil {
		return err
	}

	// 创建验证器实例
	validate := validator.New()

	// 注册翻译器
	zhTrans := zh.New()
	uni := ut.New(zhTrans, zhTrans)
	trans, _ := uni.GetTranslator("zh")
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)

	// 校验参数是否合法
	if err := validate.Struct(params); err != nil {
		return err
	}

	// 绑定成功，将参数赋值给函数参数 param
	*param = *params
	return nil
}
