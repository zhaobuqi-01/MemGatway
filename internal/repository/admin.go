package repository

import (
	"errors"
	"gateway/internal/dto"
	"gateway/pkg/utils"

	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminRepository interface {
	Getter
	Updater
	Tabler
	LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error)
}

type AdminRepo struct {
	DB *gorm.DB
}

// FindAdminByID finds an admin by their ID using GORM
func (repo *AdminRepo) Get(c *gin.Context, search *model.Admin) (*model.Admin, error) {
	return Get(c, repo.DB, "admin", search)
}

func (repo *AdminRepo) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error) {
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

func (repo *AdminRepo) Update(c *gin.Context, data *model.Admin) error {
	return Update(c, repo.DB, "admin", data)
}
