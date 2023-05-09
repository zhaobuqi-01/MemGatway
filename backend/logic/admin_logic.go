package logic

import (
	"encoding/json"
	"fmt"

	"gateway/backend/dao"
	"gateway/backend/dto"
	"gateway/backend/utils"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AdminLogic定义了管理员业务逻辑的接口
type AdminLogic interface {
	Login(c *gin.Context, params *dto.AdminLoginInput) (*dto.AdminSessionInfo, error)
	AdminLogout(c *gin.Context) error
	GetAdminInfo(c *gin.Context) (*dto.AminInfoOutput, error)
	ChangeAdminPassword(c *gin.Context, params *dto.AdminChangePwdInput) error
}

type adminLogic struct {
	db *gorm.DB
}

// NewAdminLogic创建一个新的AdminLogic实例
func NewAdminLogic(tx *gorm.DB) AdminLogic {
	return &adminLogic{
		db: tx,
	}
}

// Login验证管理员用户名和密码，并创建一个新的会话
func (s *adminLogic) Login(c *gin.Context, params *dto.AdminLoginInput) (*dto.AdminSessionInfo, error) {
	admin, err := dao.Get(c, s.db, &enity.Admin{UserName: params.UserName})
	if err != nil {
		return nil, fmt.Errorf("the username does not exist. Please try again")
	}

	if err := utils.ComparePassword(admin.Password, params.Password); err != nil {
		return nil, fmt.Errorf("incorrect password, please try again")
	}

	// 创建新的会话信息
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.ID,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}

	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		log.Error("failed to marshal session info", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal session info")
	}

	sess := sessions.Default(c)
	sess.Set(globals.AdminSessionInfoKey, string(sessBts))
	err = sess.Save()
	if err != nil {
		log.Error("failed to save session", zap.Error(err))
		return nil, fmt.Errorf("failed to save session")
	}

	return sessInfo, nil
}

// AdminLogout注销管理员会话
func (s *adminLogic) AdminLogout(c *gin.Context) error {
	sess := sessions.Default(c)
	sess.Delete(globals.AdminSessionInfoKey)
	err := sess.Save()
	if err != nil {
		log.Error("failed to save session", zap.Error(err))
		return fmt.Errorf("failed to save session")
	}
	return nil
}

// GetAdminInfo返回当前管理员的信息
func (s *adminLogic) GetAdminInfo(c *gin.Context) (*dto.AminInfoOutput, error) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(globals.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		log.Error("invalid session info", zap.Error(err))
		return nil, fmt.Errorf("invalid session info")
	}

	out := &dto.AminInfoOutput{
		ID:            adminSessionInfo.ID,
		Name:          adminSessionInfo.UserName,
		LoginTime:     adminSessionInfo.LoginTime,
		Avatar:        "https://images.unsplash.com/photo-1521747116042-5a810fda9664",
		Introduceions: "I am a super administrator",
		Roles: []string{
			"admin",
		},
	}
	return out, nil
}

// ChangeAdminPassword修改管理员密码
func (s *adminLogic) ChangeAdminPassword(c *gin.Context, params *dto.AdminChangePwdInput) error {
	sess := sessions.Default(c)
	sessInfo := sess.Get(globals.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		log.Error("invalid session info", zap.Error(err))
		return fmt.Errorf("invalid session info")
	}
	adminInfo, err := dao.Get(c, s.db, &enity.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		return fmt.Errorf("failed to get admin info")
	}

	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		log.Error("failed to generate hashed password", zap.Error(err))
		return fmt.Errorf("failed to generate hashed password")
	}

	adminInfo.Password = hashedPassword

	if err := dao.Save(c, s.db, adminInfo); err != nil {
		return fmt.Errorf("failed to change password")
	}

	return nil
}
