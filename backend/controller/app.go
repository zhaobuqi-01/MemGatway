package controller

import (
	"gateway/backend/dto"
	"gateway/backend/logic"
	"gateway/pkg/log"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App interface {
	APPList(c *gin.Context)
	APPDetail(c *gin.Context)
	APPDelete(c *gin.Context)
	APPAdd(c *gin.Context)
	APPUpdate(c *gin.Context)
}
type appController struct {
	logic.AppLogic
}

func NewAPPController(db *gorm.DB) App {
	return &appController{logic.NewAppLogic(db)}
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
// @Success 200 {object} response.Response{data=dto.APPListOutput} "success"
// @Router /app/app_list [get]
func (ac *appController) APPList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}
	list, total, err := ac.AppList(c, params)
	if err != nil {
		response.ResponseError(c, response.AppListErrCode, err)
		log.Error("Failed to fetch list", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "Get the list successfully", &dto.APPListOutput{
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
// @Success 200 {object} response.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (ac *appController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}
	app, err := ac.AppDetail(c, params)
	if err != nil {
		response.ResponseError(c, response.AppDetailErrCode, err)
		log.Error("Failed to get details", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "Get details successfully", app)
}

// APPDelete godoc
// @Summary APP删除
// @Description APP删除
// @Tags APP
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "App ID"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /app/app_delete [get]
func (ac *appController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}
	err := ac.AppDelete(c, params)
	if err != nil {
		response.ResponseError(c, response.AppDeleteErrCode, err)
		log.Error("failed to delete", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "Successfully deleted", "")
}

// APPAdd godoc
// @Summary APP添加
// @Description APP添加
// @Tags APP
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /app/app_add [post]
func (ac *appController) APPAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}
	err := ac.AppAdd(c, params)
	if err != nil {
		response.ResponseError(c, response.AppAddErrCode, err)
		log.Error("add failed", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "successfully added", "")
}

// APPUpdate godoc
// @Summary APP更新
// @Description APP更新
// @Tags APP
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /app/app_update [post]
func (ac *appController) APPUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}
	err := ac.AppUpdate(c, params)
	if err != nil {
		response.ResponseError(c, response.AppUpdateErrCode, err)
		log.Error("update failed", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "successfully updated", "")
}
