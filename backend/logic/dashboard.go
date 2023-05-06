package logic

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/dao"
	"gateway/pkg/log"
	"gateway/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardLogic interface {
	PanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error)
	ServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error)
}

type dashboardLogicImpl struct {
	db *gorm.DB
}

func NewDashboardLogic(tx *gorm.DB) *dashboardLogicImpl {
	return &dashboardLogicImpl{
		db: tx,
	}
}

// PanelGroupData 展示app数量，service数量
func (impl *dashboardLogicImpl) PanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error) {
	// 从db中分页读取基本信息
	serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%", "%")
		},
	}
	_, serviceNum, err := dao.PageList[dao.ServiceInfo](c, impl.db, serviceInfoQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get serviceNum")
	}

	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}
	_, appNum, err := dao.PageList[dao.App](c, impl.db, appQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get appNum ")
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum: serviceNum,
		AppNum:     appNum,
	}

	return out, nil
}

// ServiceStat 统计各种服务的占比
func (impl *dashboardLogicImpl) ServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error) {
	list, err := GroupByLoadType(c, impl.db)
	if err != nil {
		return nil, err
	}
	legend := []string{}
	for index, item := range list {
		name, ok := utils.LoadTypeMap[item.LoadType]
		if !ok {
			return nil, fmt.Errorf("load type not found")
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	return out, nil
}

// GroupByLoadType 按照负载类型分组
func GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
	// log记录开始查询
	log.Info("searching for group by load type", zap.String("trace_id", c.GetString("TraceID")))

	list := []dto.DashServiceStatItemOutput{}
	if err := tx.Table(dao.ServiceInfo{}.TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		log.Error("error retrieving", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, err
	}

	// log记录成功取到信息
	log.Info("group by load type was found", zap.String("trace_id", c.GetString("TraceID")))
	return list, nil
}
