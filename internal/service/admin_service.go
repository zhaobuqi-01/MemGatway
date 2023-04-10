package service

import (
	"encoding/json"
	"fmt"
	"gateway/internal/dto"
	"gateway/internal/entity"
	"gateway/internal/pkg"
	"gateway/internal/repository"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/contrib/sessions"
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
	var repo repository.Admin

	if db != nil {
		repo = repository.NewAdmin(db)
	}

	return &adminService{
		repo: repo,
	}
}

func (s *adminService) Login(c *gin.Context, params *dto.AdminLoginInput) (*dto.AdminSessionInfo, error) {
	if s.repo == nil {
		return nil, errors.New("repository is not initialized")
	}

	admin, err := s.repo.Get(c, &entity.Admin{UserName: params.UserName})
	if err != nil {
		return nil, err
	}

	saltPassword := pkg.GenSaltPassword(admin.Salt, params.Password)
	if admin.Password != saltPassword {
		return nil, errors.New("incorrect password")
	}

	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.ID,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}

	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		return nil, err
	}

	sess := sessions.Default(c)
	sess.Set(pkg.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	return sessInfo, nil
}

func (s *adminService) AdminLogout(c *gin.Context) error {
	// 业务逻辑代码
	sess := sessions.Default(c)
	sess.Delete(pkg.AdminSessionInfoKey)
	sess.Save()
	return nil
}

func (s *adminService) GetAdminInfo(c *gin.Context) (*dto.AminInfoOutput, error) {
	// 业务逻辑代码
	sess := sessions.Default(c)
	sessInfo := sess.Get(pkg.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		return nil, err
	}

	out := &dto.AminInfoOutput{
		ID:            adminSessionInfo.ID,
		Name:          adminSessionInfo.UserName,
		LoginTime:     adminSessionInfo.LoginTime,
		Avatar:        "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png",
		Introduceions: "I am a super administrator",
		Roles: []string{
			"admin",
		},
	}
	return out, nil
}

func (s *adminService) ChangeAdminPassword(c *gin.Context, params *dto.AdminChangePwdInput) error {
	if s.repo == nil {
		return errors.New("repository is not initialized")
	}
	// 业务逻辑代码
	sess := sessions.Default(c)
	sessInfo := sess.Get(pkg.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		return err
	}

	adminInfo, err := s.repo.Get(c, &entity.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		return err
	}

	saltPassword := pkg.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword

	if err := s.repo.Update(c, adminInfo); err != nil {
		return err
	}

	return nil
}
