package controller

import (
	"gateway/internal/dto"
	"gateway/internal/logic"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APPController struct {
	appLogic logic.AppLogic
}

func NewAPPController(db *gorm.DB) *APPController {
	return &APPController{appLogic: logic.NewAppLogic(db)}
}

// APPList godoc
// @Summary APP列表
// @Description APP列表
// @Tags APP
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "搜索关键字"
// @Param page_no query string true "页码"
// @Param page_size query string true "每页数量"
// @Success 200 {object} pkg.Response{data=dto.APPListOutput} "success"
// @Router /app/app_list [get]
func (ac *APPController) APPList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := params.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}
	list, total, err := ac.appLogic.AppList(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "获取app列表成功", &dto.APPListOutput{
		List:  list,
		Total: total,
	})
}

// APPDetail godoc
// @Summary APP详情
// @Description APP详情
// @Tags APP
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "App ID"
// @Success 200 {object} pkg.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (ac *APPController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}
	app, err := ac.appLogic.AppDetail(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "获取app详情成功", app)
}

// APPDelete godoc
// @Summary APP删除
// @Description APP删除
// @Tags APP
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "App ID"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /app/app_delete [get]
func (ac *APPController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}
	err := ac.appLogic.AppDelete(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "成功删除", "")
}

// APPAdd godoc
// @Summary APP添加
// @Description APP添加
// @Tags APP
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /app/app_add [post]
func (ac *APPController) APPAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}
	err := ac.appLogic.AppAdd(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "成功添加", "")
}

// APPUpdate godoc
// @Summary APP更新
// @Description APP更新
// @Tags APP
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /app/app_update [post]
func (ac *APPController) APPUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}
	err := ac.appLogic.AppUpdate(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "成功更新", "")
}

// APPStat godoc
// @Summary APP统计
// @Description APP统计
// @Tags APP
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id path string true "APP ID"
//
//	@Success 200 {object} pkg.Response{data=dto.StatisticsOutput} "success"
//
// @Router /app/app_stat [get]
func (ac *APPController) APPStat(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}
	stat, err := ac.appLogic.AppStat(c, params)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "获取app统计成功", stat)
}
