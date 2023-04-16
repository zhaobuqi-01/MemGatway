package dao

import (
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 数据库操作的通用方法

// 查询单条数据
func Get[T Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
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
func GetAll[T Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
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

func Update[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Model(data).Updates(data).Error; err != nil {
		logger.ErrorWithTraceID(c, "Error updating", zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, "Updated")
	return nil
}

// Save保存对象
func Save[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Save(data).Error; err != nil {
		logger.ErrorWithTraceID(c, "Error saving", zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, "Saved")
	return nil
}

// delete删除对象
func Delete[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Delete(data).Error; err != nil {
		logger.ErrorWithTraceID(c, "Error deleting ", zap.Error(err))
		return err
	}
	logger.InfoWithTraceID(c, "Deleted")
	return nil
}

func ListByServiceID[T Model](c *gin.Context, db *gorm.DB, serviceID int64) ([]T, int64, error) {
	var list []T
	var count int64
	query := db.Select("*")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.ErrorWithTraceID(c, "Error retrieving ", zap.Error(err))
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		logger.ErrorWithTraceID(c, "Error retrieving ", zap.Error(errCount))
		return nil, 0, err
	}
	return list, count, nil
}

// PageList 分页查询
func PageList[T Model](c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB, PageNo, PageSize int) ([]T, int64, error) {
	total := int64(0)
	list := []T{}
	offset := (PageNo - 1) * PageSize

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Limit(PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
