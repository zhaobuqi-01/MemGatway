package dao

import (
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

func (s *ServiceDetail) ServiceDetail(c *gin.Context, db *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := Get(c, db, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	log.Debug("get start")
	httpRule, err := Get(c, db, &HttpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule, err := Get(c, db, &TcpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule, err := Get(c, db, &GrpcRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl, err := Get(c, db, &AccessControl{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance, err := Get(c, db, &LoadBalance{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	log.Debug("success")
	return detail, nil
}
