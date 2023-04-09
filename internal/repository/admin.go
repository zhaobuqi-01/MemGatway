package repository

import (
	"errors"
	"gateway/internal/dto"
	"gateway/internal/model"
	"gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type adminRepository interface {
	Getter[model.Admin]
	Updater[model.Admin]
	Tabler
	LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error)
}

type adminRepo struct {
	DB *gorm.DB
}

func NewadminRepo(db *gorm.DB) *adminRepo {
	return &adminRepo{
		DB: db,
	}
}

// FindAdminByID finds an admin by their ID using GORM
func (repo *adminRepo) Get(c *gin.Context, search *model.Admin) (*model.Admin, error) {
	return Get(c, repo.DB, "admin", search)
}

func (repo *adminRepo) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error) {
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

func (repo *adminRepo) Update(c *gin.Context, data *model.Admin) error {
	return Update(c, repo.DB, "admin", data)
}
