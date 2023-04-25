package logic

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceTcpLogic interface {
	AddTCP(c *gin.Context, param *dto.ServiceAddTcpInput) error
	UpdateTCP(c *gin.Context, param *dto.ServiceUpdateTcpInput) error
}

type serviceTcpLogic struct {
	db *gorm.DB
}

func NewServiceTcpLogic(tx *gorm.DB) *serviceTcpLogic {
	return &serviceTcpLogic{db: tx}
}

func (s *serviceTcpLogic) AddTCP(c *gin.Context, params *dto.ServiceAddTcpInput) error {
	// 检查服务名是否被占用
	infoSearch := &dao.ServiceInfo{ServiceName: params.ServiceName, IsDelete: 0}
	if info, err := dao.Get(c, s.db, infoSearch); err != gorm.ErrRecordNotFound {
		if err == nil && info != nil {
			return fmt.Errorf("the TCP service name already exists, please change the service name")
		}
		return fmt.Errorf("an error occurred while querying the TCP service name")
	}

	// 检查端口是否被占用
	if _, err := dao.Get(c, s.db, &dao.TcpRule{Port: params.Port}); err == nil {
		return fmt.Errorf("the port already exists, please change the port")
	}
	if _, err := dao.Get(c, s.db, &dao.GrpcRule{Port: params.Port}); err == nil {
		return fmt.Errorf("the port already exists, please change the port")
	}

	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()
	info := &dao.ServiceInfo{
		LoadType:    pkg.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add TCP service information")
	}
	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := dao.Save(c, tx, loadBalance); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add TCP service load balancing information")
	}
	tcpRule := &dao.TcpRule{
		ServiceID: info.ID,
		Port:      params.Port,
	}
	if err := dao.Save(c, tx, tcpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add TCP service rule information")
	}
	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := dao.Save(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add TCP service permission information")
	}
	tx.Commit()
	return nil
}

func (s *serviceTcpLogic) UpdateTCP(c *gin.Context, params *dto.ServiceUpdateTcpInput) error {
	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()

	detail, err := (&dao.ServiceDetail{}).ServiceDetail(c, s.db, &dao.ServiceInfo{ID: params.ID})
	if err != nil {
		return fmt.Errorf("TCP service does not exist")
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Update(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update TCP service description")
	}

	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := dao.Update(c, tx, loadBalance); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update TCP service load balancing information")
	}

	tcpRule := &dao.TcpRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := dao.Update(c, tx, tcpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update TCP service rule information")
	}

	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := dao.Update(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update TCP service permission information")
	}

	tx.Commit()

	return nil
}
