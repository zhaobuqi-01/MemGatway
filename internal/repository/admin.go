package repository

import (
	"errors"
	"gateway/internal/dto"
	"gateway/internal/model"
	"gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Admin interface {
	Getter[model.Admin]
	Updater[model.Admin]
	LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error)
}

type admin struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) Admin {
	return &admin{
		db: db,
	}
}

// FindAdminByID finds an admin by their ID using GORM
func (repo *admin) Get(c *gin.Context, search *model.Admin) (*model.Admin, error) {
	return Get(c, repo.db, search)
}

func (repo *admin) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error) {
	adminInfo, err := repo.Get(c, &model.Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("user does not exist")
	}
	saltPassword := utils.GenSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, errors.New("wrong password, please re-enter")
	}
	return adminInfo, nil
}

func (repo *admin) Update(c *gin.Context, data *model.Admin) error {
	return Update(c, repo.db, data)
}
