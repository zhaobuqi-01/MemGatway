package controller

import (
	"gateway/internal/dto"
	"gateway/internal/repository"
	"gateway/pkg/database"
	"gateway/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type ServiceController struct{}

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
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	param := &dto.ServiceListInput{}
	if err := param.BindValParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 业务逻辑
	serviceInfo := repository.NewServiceInfo(database.GetDB())
	serviceList, total, err := serviceInfo.PageList(c, param)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	outputList := make([]dto.ServiceListItemOutput, 0)

	for _, item := range serviceList {
		outputItem := dto.ServiceListItemOutput{
			ID:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
		}
		outputList = append(outputList, outputItem)
	}

	output := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}

	middleware.ResponseSuccess(c, "", output)
}
