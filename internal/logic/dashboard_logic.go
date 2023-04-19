package logic

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardLogic interface {
	PanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error)
	ServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error)
}

type dashboardLogic struct {
	db *gorm.DB
}

func NewDashbordLogic(tx *gorm.DB) DashboardLogic {
	return &dashboardLogic{
		db: tx,
	}
}

func (d *dashboardLogic) PanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error) {
	// 从db中分页读取基本信息
	serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("service_name LIKE ? OR service_desc LIKE ?", "%", "%")
		},
	}
	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}

	_, serviceNum, err := dao.PageList[dao.ServiceInfo](c, d.db, serviceInfoQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get the total number of services")
	}
	_, appNum, err := dao.PageList[dao.App](c, d.db, appQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get the total number of services")
	}

	// counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	// if err != nil {
	// 	return
	// }
	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: 0,
		CurrentQPS:      0,
	}
	return out, nil
}

func (d *dashboardLogic) ServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error) {
	list := []dto.DashServiceStatItemOutput{}
	if err := d.db.Table((&dao.ServiceInfo{}).TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, fmt.Errorf("failed to get service statistics")
	}

	legend := []string{}
	for index, item := range list {
		name, ok := pkg.LoadTypeMap[item.LoadType]
		if !ok {
			return nil, fmt.Errorf("load_type not found")
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
