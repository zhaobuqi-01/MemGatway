package repository

import (
	"fmt"
	"gateway/internal/model"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 数据库操作的通用方法

// 查询单条数据
func Get[T model.Model](c *gin.Context, db *gorm.DB, msg string, search *T) (*T, error) {
	var out T
	result := db.Where(search).First(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.ErrorWithTraceID(c, fmt.Sprintf("%s not found ", msg), zap.Error(result.Error))
			return nil, result.Error
		}

		logger.ErrorWithTraceID(c, fmt.Sprintf("Error retrieving %s", msg), zap.Error(result.Error))
		return nil, result.Error
	}
	logger.InfoWithTraceID(c, fmt.Sprintf("Retrieved %s", msg))
	return &out, nil
}

// 查询全部
func GetAll[T model.Model](c *gin.Context, db *gorm.DB, msg string, search *T) (*T, error) {
	var out T
	result := db.Where(search).Find(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.ErrorWithTraceID(c, fmt.Sprintf("%s not found ", msg), zap.Error(result.Error))
			return nil, result.Error
		}

		logger.ErrorWithTraceID(c, fmt.Sprintf("Error retrieving %s", msg), zap.Error(result.Error))
		return nil, result.Error
	}
	logger.InfoWithTraceID(c, fmt.Sprintf("Retrieved %s", msg))
	return &out, nil
}

// create创建对象
func Create[T model.Model](c *gin.Context, db *gorm.DB, msg string, data *T) error {
	if err := db.Create(data).Error; err != nil {
		logger.ErrorWithTraceID(c, fmt.Sprintf("Error creating %s", msg), zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, fmt.Sprintf("Created %s", msg))
	return nil
}

// update更新对象
func Update[T model.Model](c *gin.Context, db *gorm.DB, msg string, data *T) error {
	if err := db.Save(data).Error; err != nil {
		logger.ErrorWithTraceID(c, fmt.Sprintf("Error updating %s", msg), zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, fmt.Sprintf("Updated %s", msg))
	return nil
}

// delete删除对象
func Delete[T model.Model](c *gin.Context, db *gorm.DB, msg string, data *T) error {
	if err := db.Delete(data).Error; err != nil {
		logger.ErrorWithTraceID(c, fmt.Sprintf("Error deleting %s", msg), zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, fmt.Sprintf("Deleted %s", msg))
	return nil
}
