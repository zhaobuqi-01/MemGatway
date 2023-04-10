package controller

import (
	"gateway/internal/dto"
	"gateway/internal/pkg"
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
		pkg.ResponseError(c, 2000, err)
		return
	}

	db := database.GetDB()

	serviceInfoService := service.NewServiceInfoService(db)
	outputList, total, err := serviceInfoService.GetServiceList(c, param)
	if err != nil {
		pkg.ResponseError(c, 2001, err)
		return
	}
	output := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}

	pkg.ResponseSuccess(c, "", output)
}
