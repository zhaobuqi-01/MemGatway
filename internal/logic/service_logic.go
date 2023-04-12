package logic

import (
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
}

type ServiceInfoLogic interface {
	Delete(c *gin.Context, param *dto.ServiceDeleteInput) error
	GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
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
