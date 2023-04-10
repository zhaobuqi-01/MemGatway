package service

import (
	"errors"
	"gateway/internal/dto"
	"gateway/internal/entity"
	"gateway/internal/pkg"
	"gateway/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminService interface {
	LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*entity.Admin, error)
}

type adminService struct {
	repo repository.Admin
}

func NewAdminService(db *gorm.DB) *adminService {
	return &adminService{
		repo: repository.NewAdmin(db),
	}
}

func (a *adminService) Get(c *gin.Context, search *entity.Admin) (*entity.Admin, error) {
	return a.repo.Get(c, search)
}

func (a *adminService) Update(c *gin.Context, data *entity.Admin) error {
	return a.repo.Update(c, data)
}

func (a *adminService) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*entity.Admin, error) {
	adminInfo, err := a.repo.Get(c, &entity.Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("user does not exist")
	}
	saltPassword := pkg.GenSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, errors.New("wrong password, please re-enter")
	}
	return adminInfo, nil
}
