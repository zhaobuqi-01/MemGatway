package controller

import (
	"gateway/internal/logic"
	"gateway/internal/pkg"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type dashboardController struct {
	dashboardLogic logic.DashboardLogic
}

func NewDashboardController(db *gorm.DB) *dashboardController {
	return &dashboardController{
		dashboardLogic: logic.NewDashbordLogic(db),
	}
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]

func (dc *dashboardController) PanelGroupData(c *gin.Context) {
	out, err := dc.dashboardLogic.PanelGroupData(c)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		log.Error("PanelGroupData failed", zap.Error(err))
		return
	}
	pkg.ResponseSuccess(c, "", out)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (dc *dashboardController) ServiceStat(c *gin.Context) {
	out, err := dc.dashboardLogic.ServiceStat(c)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		log.Error("ServiceStat failed", zap.Error(err))
		return
	}
	pkg.ResponseSuccess(c, "", out)
}
