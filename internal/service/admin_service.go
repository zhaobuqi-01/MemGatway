package service

import (
	"errors"
	"gateway/internal/dto"
	"gateway/internal/model"
	"gateway/internal/pkg"
	"gateway/internal/repository"

	"github.com/gin-gonic/gin"
)

type AdminService interface {
	LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error)
}

type adminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) AdminService {
	return &adminService{
		repo: repo,
	}
}
func (a *adminService) LoginCheck(c *gin.Context, param *dto.AdminLoginInput) (*model.Admin, error) {
	adminInfo, err := a.repo.Get(c, &model.Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("user does not exist")
	}
	saltPassword := pkg.GenSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, errors.New("wrong password, please re-enter")
	}
	return adminInfo, nil
}
