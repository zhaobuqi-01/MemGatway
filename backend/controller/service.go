package controller

import (
	"gateway/backend/dto"
	"gateway/backend/logic"
	"gateway/pkg/log"
	"gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service interface {
	ServiceList(c *gin.Context)
	ServiceDelete(c *gin.Context)
	ServiceDetail(c *gin.Context)
	ServiceAddHttp(c *gin.Context)
	ServiceUpdateHttp(c *gin.Context)
	ServiceAddTcp(c *gin.Context)
	ServiceUpdateTcp(c *gin.Context)
	ServiceAddGrpc(c *gin.Context)
	ServiceUpdateGrpc(c *gin.Context)
	ServiceStat(c *gin.Context)
}
type serviceController struct {
	logic.ServiceLogic
}

func NewServiceController() *serviceController {
	return &serviceController{
		logic.NewServiceLogic(),
	}
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags Service
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query  string false "关键词"
// @Param page_no query  int true "页码"
// @Param page_size query  int true "每页条数"
// @Success 200 {object} response.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *serviceController) ServiceList(c *gin.Context) {
	param := &dto.ServiceListInput{}
	if err := param.BindValParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	outputList, total, err := s.GetServiceList(c, param)
	if err != nil {
		response.ResponseError(c, response.ServiceListErrCode, err)
		log.Error("Failed to fetch list", zap.Error(err))
		return
	}

	output := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}

	response.ResponseSuccess(c, "", output)
}

// ServiceAdd godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags Service
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query  int true "服务id"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_delete [get]
func (s *serviceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := s.DeleteService(c, params)
	if err != nil {
		response.ResponseError(c, response.ServiceDeleteErrCode, err)
		log.Error("Failed to delete service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "", "delete success")
}

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags Service
// ID /service/service_detail
// Accept json
// Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} response.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (s *serviceController) ServiceDetail(c *gin.Context) {

	param := &dto.ServiceDeleteInput{}

	if err := param.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	output, err := s.GetServiceDetail(c, param)

	if err != nil {
		response.ResponseError(c, response.ServiceDetailErrCode, err)
		log.Error("Failed to get service detail", zap.Error(err))
		return

	}

	response.ResponseSuccess(c, "", output)

}

// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags Service
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (s *serviceController) ServiceAddHttp(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := s.AddHTTP(c, params)
	if err != nil {
		response.ResponseError(c, response.AddHttpServiceErrCode, err)
		log.Error("Failed to add http service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "add httpService success", nil)
}

// ServiceUpadteHTTP godoc
// @Summary 更新HTTP服务
// @Description 更新HTTP服务
// @Tags Service
// @ID /service/service_update_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_update_http [post]
func (s *serviceController) ServiceUpdateHttp(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	if err := params.BindValParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}
	log.Debug("httpService params", zap.Any("params", params))
	err := s.UpdateHTTP(c, params)
	if err != nil {
		response.ResponseError(c, response.UpdateHttpServiceErrCode, err)
		log.Error("Failed to update http service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "update httpService success", nil)
}

// ServiceAddTcp godoc
// @Summary 添加TCP服务
// @Description 添加TCP服务
// @Tags Service
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (s *serviceController) ServiceAddTcp(c *gin.Context) {
	params := &dto.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := s.AddTCP(c, params)
	if err != nil {
		response.ResponseError(c, response.AddTCPServiceErrCode, err)
		log.Error("Failed to add tcp service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "add tcpService success", nil)
}

// ServiceUpdateTcp godoc
// @Summary 更新TCP服务
// @Description 更新TCP服务
// @Tags Service
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (s *serviceController) ServiceUpdateTcp(c *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := s.UpdateTCP(c, params)
	if err != nil {
		response.ResponseError(c, response.UpdateTCPServiceErrCode, err)
		log.Error("Failed to update tcp service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "update tcpService success", nil)
}

// ServiceAddGrpc godoc
// @Summary 添加GRPC服务
// @Description 添加GRPC服务
// @Tags Service
// @ID /service/service_add_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_add_grpc [post]
func (s *serviceController) ServiceAddGrpc(c *gin.Context) {
	params := &dto.ServiceAddGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := s.AddGrpc(c, params)
	if err != nil {
		response.ResponseError(c, response.AddGRPCServiceErrCode, err)
		log.Error("Failed to add grpc service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "add grpcService success", nil)
}

// ServiceUpdateGrpc godoc
// @Summary 更新GRPC服务
// @Description 更新GRPC服务
// @Tags Service
// @ID /service/service_update_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} response.Response{data=string} "success"
// @Router /service/service_update_grpc [post]
func (s *serviceController) ServiceUpdateGrpc(c *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	err := s.UpdateGrpc(c, params)
	if err != nil {
		response.ResponseError(c, response.UpdateGRPCServiceErrCode, err)
		log.Error("Failed to update grpc service", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, "update grpcService success", nil)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags Service
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} response.Response{data=dto.ServiceStatOutput}
// @Router /service/service_stat [get]
func (s *serviceController) ServiceStat(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		response.ResponseError(c, response.ParamBindingErrCode, err)
		return
	}

	output, err := s.GetServiceStat(c, params)
	if err != nil {
		response.ResponseError(c, response.CommErrCode, err)
		log.Error("Failed to get service stat", zap.Error(err))
		return
	}

	response.ResponseSuccess(c, "", output)
}
