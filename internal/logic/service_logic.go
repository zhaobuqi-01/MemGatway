package logic

import (
	"gateway/internal/dao"
	"gateway/internal/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceLogic interface {
	ServiceInfoLogic
	ServiceTcpLogic
	ServiceHttpLogic
	ServiceGrpcLogic
}

type ServiceInfoLogic interface {
	ServiceDelete(c *gin.Context, param *dto.ServiceDeleteInput) error
	GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
	GetServiceDetail(c *gin.Context, param *dto.ServiceDeleteInput) (*dao.ServiceDetail, error)
	GetServiceStat(c *gin.Context, param *dto.ServiceDeleteInput) (*dto.ServiceStatOutput, error)
}

type ServiceTcpLogic interface {
	AddTCP(c *gin.Context, param *dto.ServiceAddTcpInput) error
	UpdateTCP(c *gin.Context, param *dto.ServiceUpdateTcpInput) error
}

type ServiceHttpLogic interface {
	AddHTTP(c *gin.Context, param *dto.ServiceAddHTTPInput) error
	UpdateHTTP(c *gin.Context, param *dto.ServiceUpdateHTTPInput) error
}

type ServiceGrpcLogic interface {
	AddGrpc(c *gin.Context, param *dto.ServiceAddGrpcInput) error
	UpdateGrpc(c *gin.Context, param *dto.ServiceUpdateGrpcInput) error
}
type serviceLogic struct {
	ServiceInfoLogic
	ServiceHttpLogic
	ServiceTcpLogic
	ServiceGrpcLogic
}

func NewServiceLogic(db *gorm.DB) ServiceLogic {
	return &serviceLogic{
		ServiceInfoLogic: NewServiceInfoLogic(db),
		ServiceHttpLogic: NewServiceHttpLogic(db),
		ServiceTcpLogic:  NewServiceTcpLogic(db),
		ServiceGrpcLogic: NewServiceGrpcLogic(db),
	}
}
