package repository

import (
	"errors"
	"gateway/internal/dto"
	"gateway/pkg/utils"

	"gateway/internal/model"

	"github.com/gin-gonic/gin"
)

type AdminRepository interface {
	Getter
	Updater
	Tabler
	LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error)
}

type AdminRepo struct{}

// 明确表名为admin而不是默认表名admins
func (repo *AdminRepo) TableName() string {
	return "gateway_admin"
}

// FindAdminByID finds an admin by their ID using GORM
func (repo *AdminRepo) GetAll(c *gin.Context, search *model.Admin) (*model.Admin, error) {
	return GetAll(c, "admin", search)
}

func (repo *AdminRepo) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error) {
	adminInfo, err := repo.GetAll(c, &model.Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("User does not exist")
	}
	saltPassword := utils.GenSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, errors.New("Wrong password, please re-enter")
	}
	return adminInfo, nil
}
