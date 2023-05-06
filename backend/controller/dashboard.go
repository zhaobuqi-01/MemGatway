package controller

import (
	"gateway/backend/logic"
	"gateway/pkg/log"
	"gateway/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type dashboardController struct {
	logic logic.DashboardLogic
}

func NewDashboardController(db *gorm.DB) *dashboardController {
	return &dashboardController{
		logic: logic.NewDashboardLogic(db),
	}
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags Dashboard
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (dc *dashboardController) PanelGroupData(c *gin.Context) {
	out, err := dc.logic.PanelGroupData(c)
	if err != nil {
		utils.ResponseError(c, utils.PanelGroupDataErrCode, err)
		log.Error("failed to get data", zap.Error(err))
		return
	}
	utils.ResponseSuccess(c, "get data successfully", out)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags Dashboard
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (dc *dashboardController) ServiceStat(c *gin.Context) {
	out, err := dc.logic.ServiceStat(c)
	if err != nil {
		utils.ResponseError(c, utils.ServiceStatErrCode, err)
		log.Error("failed to get serviceStat", zap.Error(err))
		return
	}
	utils.ResponseSuccess(c, "get serviceStat successfully", out)
}
