package controller

import (
	"gateway/backend/logic"
	"gateway/pkg/log"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Dashboard interface {
	PanelGroupData(c *gin.Context)
	ServiceStat(c *gin.Context)
}

type dashboardController struct {
	logic.DashboardLogic
}

func NewDashboardController() *dashboardController {
	return &dashboardController{logic.NewDashboardLogic()}
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags Dashboard
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (dc *dashboardController) PanelGroupData(c *gin.Context) {
	out, err := dc.GetPanelGroupData(c)
	if err != nil {
		response.ResponseError(c, response.PanelGroupDataErrCode, err)
		log.Error("failed to get data", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "get data successfully", out)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags Dashboard
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (dc *dashboardController) ServiceStat(c *gin.Context) {
	out, err := dc.GetServiceStat(c)
	if err != nil {
		response.ResponseError(c, response.ServiceStatErrCode, err)
		log.Error("failed to get serviceStat", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "get serviceStat successfully", out)
}

// FlowStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (service *dashboardController) FlowStat(c *gin.Context) {
	out, err := service.GetFlowStat(c)
	if err != nil {
		response.ResponseError(c, response.CommErrCode, err)
		log.Error("failed to get serviceStat", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "get serviceStat successfully", out)
}
