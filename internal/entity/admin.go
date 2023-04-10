package entity

import "time"

// Admin表对应的实体类
type Admin struct {
	ID       int       `json:"id" gorm:"primary_key" description:"主键"`
	UserName string    `json:"user_name" gorm:"column:user_name" description:"用户名"`
	Salt     string    `json:"salt" gorm:"column:salt" description:"盐值"`
	Password string    `json:"password" gorm:"column:password" description:"密码"`
	UpdateAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (Admin) TableName() string {
	return "gateway_admin"
}
