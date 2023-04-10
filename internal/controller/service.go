package controller

import (
	"gateway/internal/dto"
	"gateway/internal/middleware"
	"gateway/internal/repository"
	"gateway/internal/service"
	"gateway/pkg/database"

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
func (s *ServiceController) ServiceList(c *gin.Context) {
	param := &dto.ServiceListInput{}
	if err := param.BindValParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db := database.GetDB()
	serviceInfoRepo := repository.NewServiceInfoRepo(db)

	// 从db中分页读取基本信息
	list, total, err := serviceInfoRepo.PageList(c, param)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 格式化输出信息
	outputList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := serviceInfoRepo.ServiceDetail(c, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			return
		}

		// 根据服务类型和规则生成服务地址
		serviceInfoService := service.NewServiceInfoService()
		serviceAddr, err := serviceInfoService.GetServiceAddress(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}

		// 获取IP列表
		ipList := serviceInfoService.GetIPList(c, serviceDetail.LoadBalance)

		outputItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(ipList),
		}
		outputList = append(outputList, outputItem)
	}

	output := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}

	middleware.ResponseSuccess(c, "", output)
}
