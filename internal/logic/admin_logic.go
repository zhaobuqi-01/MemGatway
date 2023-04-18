package logic

import (
	"encoding/json"
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminLogic interface {
	Login(c *gin.Context, params *dto.AdminLoginInput) (*dto.AdminSessionInfo, error)
	AdminLogout(c *gin.Context) error
	GetAdminInfo(c *gin.Context) (*dto.AminInfoOutput, error)
	ChangeAdminPassword(c *gin.Context, params *dto.AdminChangePwdInput) error
}

type adminLogic struct {
	db *gorm.DB
}

func NewAdminLogic(tx *gorm.DB) *adminLogic {
	return &adminLogic{
		db: tx,
	}
}

func (s *adminLogic) Login(c *gin.Context, params *dto.AdminLoginInput) (*dto.AdminSessionInfo, error) {
	admin, err := dao.Get(c, s.db, &dao.Admin{UserName: params.UserName})
	if err != nil {
		return nil, fmt.Errorf("用户名不存在，请重新输入！")
	}

	if err := pkg.ComparePassword(admin.Password, params.Password); err != nil {
		return nil, fmt.Errorf("密码错误，请重新输入！")
	}

	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.ID,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}

	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		return nil, fmt.Errorf("session info marshal failed")
	}

	sess := sessions.Default(c)
	sess.Set(pkg.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	return sessInfo, nil
}

func (s *adminLogic) AdminLogout(c *gin.Context) error {
	// 业务逻辑代码
	sess := sessions.Default(c)
	sess.Delete(pkg.AdminSessionInfoKey)
	sess.Save()
	return nil
}

func (s *adminLogic) GetAdminInfo(c *gin.Context) (*dto.AminInfoOutput, error) {
	// 业务逻辑代码
	sess := sessions.Default(c)
	sessInfo := sess.Get(pkg.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		return nil, fmt.Errorf("session info is not valid")
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

func (s *adminLogic) ChangeAdminPassword(c *gin.Context, params *dto.AdminChangePwdInput) error {
	// 业务逻辑代码
	sess := sessions.Default(c)
	sessInfo := sess.Get(pkg.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		return fmt.Errorf("session info is not valid")
	}

	adminInfo, err := dao.Get(c, s.db, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		return fmt.Errorf("admin.Get: %w", err)
	}

	hashedPassword, err := pkg.HashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("GenSaltPassword: %w", err)
	}

	adminInfo.Password = hashedPassword

	if err := dao.Update(c, s.db, adminInfo); err != nil {
		return fmt.Errorf("admin.Update: %w", err)
	}

	return nil
}
