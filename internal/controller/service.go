package controller

import (
	"gateway/internal/dto"
	"gateway/internal/logic"
	"gateway/internal/pkg"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceController struct {
	logic logic.ServiceLogic
}

func NewServiceController(db *gorm.DB) *ServiceController {
	return &ServiceController{
		logic: logic.NewServiceLogic(db),
	}
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务接口
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query  string false "关键词"
// @Param page_no query  int true "页码"
// @Param page_size query  int true "每页条数"
// @Success 200 {object} pkg.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *ServiceController) ServiceList(c *gin.Context) {
	param := &dto.ServiceListInput{}
	if err := param.BindValParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}

	outputList, total, err := s.logic.GetServiceList(c, param)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}

	output := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}

	pkg.ResponseSuccess(c, "", output)
}

// ServiceAdd godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务接口
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query  int true "服务id"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /service/service_delete [get]
func (s *ServiceController) ServiceDelete(c *gin.Context) {
	param := &dto.ServiceDeleteInput{}
	if err := param.BindValidParam(c); err != nil {
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}

	err := s.logic.Delete(c, param)
	if err != nil {
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}

	pkg.ResponseSuccess(c, "", "delete success")
}

// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务接口
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} pkg.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (s *ServiceController) ServiceAddHttp(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValParam(c); err != nil {
		logger.ErrorWithTraceID(c, "parameter binding error")
		pkg.ResponseError(c, pkg.ParamBindingErrCode, err)
		return
	}

	err := s.logic.AddHTTP(c, params)
	if err != nil {
		logger.ErrorWithTraceID(c, "service add http error")
		pkg.ResponseError(c, pkg.InternalErrorCode, err)
		return
	}
	pkg.ResponseSuccess(c, "add httpService success", nil)
}
