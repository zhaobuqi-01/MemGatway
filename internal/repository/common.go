package repository

import (
	"fmt"
	mysql "gateway/pkg/database/mysql"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 泛型函数用于实现通用的数据库操作
// go没有泛型方法所以不能以方法的形式实现
// 但是可以通过泛型函数实现

// 数据库操作的通用方法

type Model interface {
	Admin | User
}

// Find[T model.Model]是泛型函数，T是泛型参数，model.Model是泛型约束
// Find检索全部对象
func GetAll[T Model](c *gin.Context, describe string, search *T) (*T, error) {
	logger.InfoWithTraceID(c, "开始从数据库获取 %s 信息", describe)
	db := mysql.GetDB()

	var out T
	result := db.Where(search).Find(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.ErrorWithTraceID(c, fmt.Sprintf("%s not found ", describe), zap.Error(result.Error))
			return nil, result.Error
		}

		logger.ErrorWithTraceID(c, fmt.Sprintf("Error retrieving %s", describe), zap.Error(result.Error))
		return nil, result.Error
	}
	logger.InfoWithTraceID(c, "从数据库获取 %s 信息成功", describe)
	return &out, nil
}
