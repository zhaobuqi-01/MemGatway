package logic

import (
	"fmt"
	"gateway/backend/dao"
	"gateway/backend/dto"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardLogic interface {
	GetPanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error)
	GetServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error)
	GetFlowStat(c *gin.Context) (*dto.ServiceStatOutput, error)
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
func (impl *dashboardLogicImpl) GetPanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error) {
	// 从db中分页读取基本信息
	serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%", "%")
		},
	}
	_, serviceNum, err := dao.PageList[enity.ServiceInfo](c, impl.db, serviceInfoQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get serviceNum")
	}

	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}
	_, appNum, err := dao.PageList[enity.App](c, impl.db, appQueryConditions, 1, 1)
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
func (impl *dashboardLogicImpl) GetServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error) {
	list, err := dao.GetLoadTypeByGroup(c, impl.db)
	if err != nil {
		return nil, err
	}
	legend := []string{}
	for index, item := range list {
		name, ok := globals.LoadTypeMap[item.LoadType]
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

// flow_stat流量统计
func (impl *dashboardLogicImpl) GetFlowStat(c *gin.Context) (*dto.ServiceStatOutput, error) {
	counter, err := globals.FlowCounter.GetCounter(globals.FlowTotal)
	if err != nil {
		log.Error("failed to get flow counter ", zap.Error(err))
		return nil, err
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, time.UTC)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, time.UTC)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	out := &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}
	return out, nil
}
