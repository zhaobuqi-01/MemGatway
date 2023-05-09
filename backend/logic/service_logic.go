package logic

import (
	"gorm.io/gorm"
)

type ServiceLogic interface {
	ServiceInfoLogic
	TcpServiceLogic
	HttpServiceLogic
	GrpcServiceLogic
}

type serviceLogic struct {
	ServiceInfoLogic
	HttpServiceLogic
	TcpServiceLogic
	GrpcServiceLogic
}

func NewServiceLogic(db *gorm.DB) ServiceLogic {
	return &serviceLogic{
		ServiceInfoLogic: NewServiceInfoLogic(db),
		HttpServiceLogic: NewHttpServiceLogic(db),
		TcpServiceLogic:  NewTcpServiceLogic(db),
		GrpcServiceLogic: NewGrpcServiceLogic(db),
	}
}
