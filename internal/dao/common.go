package dao

import (
	"fmt"
	"gateway/pkg/log"

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
			log.Error(fmt.Sprintf(" %v not found ", search), zap.Error(result.Error), zap.String("trace_id", c.GetString("TraceID")))
			return nil, result.Error
		}

		log.Error(fmt.Sprintf("error retrieving :%v ", search), zap.Error(result.Error), zap.String("trace_id", c.GetString("TraceID")))
		return nil, result.Error
	}
	log.Info(fmt.Sprintf(" %v was found", search), zap.String("trace_id", c.GetString("TraceID")))
	return &out, nil
}

// update更新对象
func Update[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Model(data).Updates(data).Error; err != nil {
		log.Error(fmt.Sprintf("error updating : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}
	log.Info(fmt.Sprintf("%v updated", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// Save保存对象
func Save[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Save(data).Error; err != nil {
		log.Error(fmt.Sprintf("error saving : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}
	log.Info(fmt.Sprintf("%v Saved", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// delete删除对象
func Delete[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	if err := db.Delete(data).Error; err != nil {
		log.Error(fmt.Sprintf("error deleting : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}
	log.Info(fmt.Sprintf("%v deleted", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

func ListByServiceID[T Model](c *gin.Context, db *gorm.DB, serviceID int64) ([]T, int64, error) {
	var list []T
	var count int64
	query := db.Select("*")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", serviceID), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		log.Error(fmt.Sprintf("error retrieving :%v ", serviceID), zap.Error(errCount), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	log.Info(fmt.Sprintf("%v was found", serviceID), zap.String("trace_id", c.GetString("TraceID")))
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
		log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	query.Limit(PageSize).Offset(offset).Count(&total)
	log.Info(fmt.Sprintf("%v was found", query), zap.String("trace_id", c.GetString("TraceID")))
	return list, total, nil
}
