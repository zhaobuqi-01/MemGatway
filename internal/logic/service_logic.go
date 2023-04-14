package logic

import (
	"gateway/internal/dao"
	"gateway/internal/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceLogic interface {
	ServiceInfoLogic
	ServiceHttpLogic
}

type ServiceHttpLogic interface {
	AddHTTP(c *gin.Context, param *dto.ServiceAddHTTPInput) error
	UpdateHTTP(c *gin.Context, param *dto.ServiceUpdateHTTPInput) error
}

type ServiceInfoLogic interface {
	Delete(c *gin.Context, param *dto.ServiceDeleteInput) error
	GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
	GetServiceDetail(c *gin.Context, param *dto.ServiceDeleteInput) (*dao.ServiceDetail, error)
}

type serviceLogic struct {
	ServiceInfoLogic
	ServiceHttpLogic
}

func NewServiceLogic(db *gorm.DB) ServiceLogic {
	return &serviceLogic{
		ServiceInfoLogic: NewServiceInfoLogic(db),
		ServiceHttpLogic: NewServiceHttpLogic(db),
	}
}
