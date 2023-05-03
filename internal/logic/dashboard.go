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

type dashboardLogicImpl struct {
	db *gorm.DB
}

func NewDashboardLogic(tx *gorm.DB) *dashboardLogicImpl {
	return &dashboardLogicImpl{
		db: tx,
	}
}

func (impl *dashboardLogicImpl) PanelGroupData(c *gin.Context) (*dto.PanelGroupDataOutput, error) {
	// 从db中分页读取基本信息
	serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%"+""+"%", "%"+""+"%")
		},
	}
	_, serviceNum, err := dao.PageList[dao.ServiceInfo](c, impl.db, serviceInfoQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get serviceNum")
	}

	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%"+""+"%", "%"+""+"%")
		},
	}
	_, appNum, err := dao.PageList[dao.App](c, impl.db, appQueryConditions, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get appNum ")
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: 0,
		CurrentQPS:      0,
	}

	return out, nil
}

func (impl *dashboardLogicImpl) ServiceStat(c *gin.Context) (*dto.DashServiceStatOutput, error) {
	list, err := dao.GroupByLoadType(c, impl.db)
	if err != nil {
		return nil, err
	}
	legend := []string{}
	for index, item := range list {
		name, ok := pkg.LoadTypeMap[item.LoadType]
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
