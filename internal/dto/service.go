package dto

import (
	"gateway/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词"  validate:""`                                //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页码" example:"1" validate:"required"`        //页码
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` //每页条数
}

func (params *ServiceListInput) BindValParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}

type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id" comment:"服务ID"  validate:""`                     //服务ID
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称"  validate:""` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述"  validate:""` //服务描述
	LoadType    int    `json:"load_type" form:"load_type" comment:"负载类型"  validate:""`       //负载类型
	ServiceAddr string `json:"service_addr" form:"service_addr" comment:"服务地址"  validate:""` //服务地址
	Qps         int    `json:"qps" form:"qps" comment:"QPS"  validate:""`                    //QPS
	Qpd         int    `json:"qpd" form:"qpd" comment:"QPD"  validate:""`                    //QPD
	TotalNode   int    `json:"total_node" form:"total_node" comment:"节点总数"  validate:""`     //节点总数

}

type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数"  validate:""` //总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表"  validate:""`   //列表
}
