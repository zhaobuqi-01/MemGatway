package logic

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/dao"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"
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
	service dao.AllGetter[enity.ServiceInfo]
	getData dao.LoadTypeGrouper[enity.ServiceInfo]
	app     dao.AllGetter[enity.App]
	db      *gorm.DB
}

func NewDashboardLogic() *dashboardLogicImpl {
	return &dashboardLogicImpl{
		dao.New[enity.ServiceInfo](),
		dao.New[enity.ServiceInfo](),
		dao.New[enity.App](),
		mysql.GetDB(),
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

	log.Debug("start to get serviceNum")
	serviceList, err := impl.service.GetAll(c, impl.db, serviceInfoQueryConditions)
	if err != nil {
		return nil, fmt.Errorf("failed to get serviceNum")
	}
	log.Debug("end to get serviceNum", zap.Int("serviceNum", len(serviceList)))

	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}

	log.Debug("start to get appNum")
	appList, err := impl.app.GetAll(c, impl.db, appQueryConditions)
	if err != nil {
		return nil, fmt.Errorf("failed to get appNum ")
	}
	log.Debug("end to get appNum", zap.Int("appNum", len(appList)))

	counter, err := globals.FlowCounter.GetCounter(globals.FlowTotal)
	if err != nil {
		return nil, fmt.Errorf("get flow counter failed")
	}
	out := &dto.PanelGroupDataOutput{
		ServiceNum:      int64(len(serviceList)),
		AppNum:          int64(len(appList)),
		TodayRequestNum: counter.QPD,
		CurrentQPS:      counter.QPS,
	}

	return out, nil
}

// ServiceStat 统计各种服务的占比
func (impl *dashboardLogicImpl) GetServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error) {
	list, err := impl.getData.GetLoadTypeByGroup(c, impl.db)
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

	log.Debug("get flow stat successfully", zap.Any("data", out))
	return out, nil
}
