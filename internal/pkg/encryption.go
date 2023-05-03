package pkg

import (
	"gateway/pkg/log"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// GenSaltPassword 使用bcrypt替换
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("GenSaltPassword failed", zap.Error(err))
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword 对比密码
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
