package repository

import (
	"errors"
	"gateway/internal/dto"
	"gateway/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminRepository interface {
	Getter
	Updater
	Tabler
}

// Admin表对应的实体类
type Admin struct {
	Id       int64     `json:"id" gorm:"primary_key" description:"主键"`
	UserName string    `json:"user_name" gorm:"column:user_name" description:"用户名"`
	Salt     string    `json:"salt" gorm:"column:salt" description:"盐值"`
	Password string    `json:"password" gorm:"column:password" description:"密码"`
	UserId   int64     `json:"user_id" gorm:"column:user_id" description:"用户id"`
	UpateAt  time.Time `json:"upate_at" gorm:"column:upate_at" description:"更新时间"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

// 明确表名为admin而不是默认表名admins
func (admin *Admin) TableName() string {
	return "gateway_admin"
}

// FindAdminByID finds an admin by their ID using GORM
func (admin *Admin) GetAll(c *gin.Context, search *Admin) (*Admin, error) {
	return GetAll(c, "admin", search)
}

func (admin *Admin) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := admin.GetAll(c, &Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("User does not exist")
	}
	saltPassword := utils.GenSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, errors.New("Wrong password, please re-enter")
	}
	return adminInfo, nil
}
