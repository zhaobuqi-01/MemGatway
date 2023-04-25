package logic

import (
	"gorm.io/gorm"
)

type ServiceLogic interface {
	ServiceInfoLogic
	ServiceTcpLogic
	ServiceHttpLogic
	ServiceGrpcLogic
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
