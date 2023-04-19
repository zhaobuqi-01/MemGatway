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

type servcieGrpcLogic struct {
	db *gorm.DB
}

func NewServiceGrpcLogic(tx *gorm.DB) *servcieGrpcLogic {
	return &servcieGrpcLogic{db: tx}
}

func (s *servcieGrpcLogic) AddGrpc(c *gin.Context, params *dto.ServiceAddGrpcInput) error {
	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if info, err := dao.Get(c, s.db, infoSearch); err != gorm.ErrRecordNotFound {
		if err == nil && info != nil {
			return fmt.Errorf("the GRPC service name already exists, please change the service name")
		}
		return fmt.Errorf("an error occurred while querying the service name")
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := dao.Get(c, s.db, tcpRuleSearch); err == nil {
		return fmt.Errorf("the port already exists, please change the port")
	}

	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := dao.Get(c, s.db, grpcRuleSearch); err == nil {
		return fmt.Errorf("the port already exists, please change the port")
	}

	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the ip list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()

	info := &dao.ServiceInfo{
		LoadType:    pkg.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add GRPC service information")
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
		return fmt.Errorf("failed to add GRPC service load balancer")
	}

	grpcRule := &dao.GrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := dao.Save(c, tx, grpcRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add GRPC service rule")
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
		return fmt.Errorf("failed to add GRPC service permission")
	}

	tx.Commit()

	return nil
}

func (s *servcieGrpcLogic) UpdateGrpc(c *gin.Context, params *dto.ServiceUpdateGrpcInput) error {
	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the ip list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := (&dao.ServiceDetail{}).ServiceDetail(c, s.db, service)
	if err != nil {
		return fmt.Errorf("gRPC service does not exist")
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Update(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update GRPC service information")
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
		return fmt.Errorf("failed to update GRPC service load balancing")
	}

	grpcRule := &dao.GrpcRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	// grpcRule.Port = params.Port
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := dao.Update(c, tx, grpcRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update GRPC service rules")
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
		return fmt.Errorf("failed to update GRPC service permissions")
	}

	tx.Commit()

	return nil
}
