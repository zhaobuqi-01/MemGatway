package repository

import (
	"gateway/internal/model"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 数据库操作的通用方法

// 查询单条数据
func Get[T model.Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
	var out T
	result := db.Where(search).First(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.ErrorWithTraceID(c, "not found", zap.Error(result.Error))
			return nil, result.Error
		}

		logger.ErrorWithTraceID(c, "Error retrieving ", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.InfoWithTraceID(c, "Retrieved ")
	return &out, nil
}

// 查询全部
func GetAll[T model.Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
	var out T
	result := db.Where(search).Find(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.ErrorWithTraceID(c, "not found", zap.Error(result.Error))
		}

		logger.ErrorWithTraceID(c, "Error retrieving ", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.InfoWithTraceID(c, "Retrieved ")
	return &out, nil
}

// create创建对象
func Create[T model.Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Create(data).Error; err != nil {
		logger.ErrorWithTraceID(c, "Error creating ", zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, "Created")
	return nil
}

// update更新对象
func Update[T model.Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Save(data).Error; err != nil {
		logger.ErrorWithTraceID(c, "Error updating", zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, "Updated")
	return nil
}

// delete删除对象
func Delete[T model.Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Delete(data).Error; err != nil {
		logger.ErrorWithTraceID(c, "Error deleting ", zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, "Deleted")
	return nil
}
