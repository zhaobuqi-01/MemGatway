package logic

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

func NewServiceLogic() *serviceLogic {
	return &serviceLogic{
		ServiceInfoLogic: NewServiceInfoLogic(),
		HttpServiceLogic: NewHttpServiceLogic(),
		TcpServiceLogic:  NewTcpServiceLogic(),
		GrpcServiceLogic: NewGrpcServiceLogic(),
	}
}
