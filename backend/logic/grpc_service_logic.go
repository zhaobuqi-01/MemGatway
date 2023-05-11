package logic

import (
	"fmt"
	"gateway/backend/dao"
	"gateway/backend/dto"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"

	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GrpcServiceLogic interface {
	AddGrpc(c *gin.Context, param *dto.ServiceAddGrpcInput) error
	UpdateGrpc(c *gin.Context, param *dto.ServiceUpdateGrpcInput) error
}

// grpcServiceLogic 结构体
type grpcServiceLogic struct {
	db *gorm.DB
}

// NewGrpcServiceLogic 构造函数
func NewGrpcServiceLogic(tx *gorm.DB) *grpcServiceLogic {
	return &grpcServiceLogic{
		db: tx,
	}
}

// AddGrpc 添加 GRPC 服务
func (s *grpcServiceLogic) AddGrpc(c *gin.Context, params *dto.ServiceAddGrpcInput) error {
	// 验证 service_name 是否被占用
	infoSearch := &enity.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if info, err := dao.Get(c, s.db, infoSearch); err != gorm.ErrRecordNotFound {
		if err == nil && info != nil {
			return fmt.Errorf("GRPC service name already exists, please change the service name")
		}
		return fmt.Errorf("error occurred while querying the service name")
	}

	// 验证端口是否被占用
	tcpRuleSearch := &enity.TcpRule{
		Port: params.Port,
	}
	if _, err := dao.Get(c, s.db, tcpRuleSearch); err == nil {
		return fmt.Errorf("port already exists, please change the port")
	}

	grpcRuleSearch := &enity.GrpcRule{
		Port: params.Port,
	}
	if _, err := dao.Get(c, s.db, grpcRuleSearch); err == nil {
		return fmt.Errorf("port already exists, please change the port")
	}

	// 检查 IP 列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("ip list is inconsistent with the number of weight lists")
	}

	// 开始事务
	tx := s.db.Begin()

	// 保存服务信息
	info := &enity.ServiceInfo{
		LoadType:    globals.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add GRPC service information")
	}

	// 保存负载均衡信息
	loadBalance := &enity.LoadBalance{
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

	// 保存 GRPC 规则
	grpcRule := &enity.GrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := dao.Save(c, tx, grpcRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add GRPC service rule")
	}

	// 保存访问控制信息
	accessControl := &enity.AccessControl{
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

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:        "service",
		Payload:     params.ServiceName,
		ServiceType: globals.LoadTypeGRPC,
		Operation:   globals.DataInsert,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// UpdateGrpc 更新 GRPC 服务
func (s *grpcServiceLogic) UpdateGrpc(c *gin.Context, params *dto.ServiceUpdateGrpcInput) error {
	// 检查 IP 列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("ip list is inconsistent with the number of weight lists")
	}
	// 开始事务
	tx := s.db.Begin()

	// 获取服务详情
	service := &enity.ServiceInfo{
		ID: params.ID,
	}
	detail, err := dao.GetServiceDetail(c, s.db, service)
	if err != nil {
		return fmt.Errorf("gRPC service does not exist")
	}

	// 更新服务信息
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save GRPC service information")
	}

	// 更新负载均衡信息
	loadBalance := &enity.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := dao.Save(c, tx, loadBalance); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save GRPC service load balancing")
	}

	// 更新 GRPC 规则
	grpcRule := &enity.GrpcRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := dao.Save(c, tx, grpcRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save GRPC service rules")
	}

	// 更新访问控制信息
	accessControl := &enity.AccessControl{}
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
	if err := dao.Save(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save GRPC service permissions")
	}

	// 提交事务
	tx.Commit()

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:        "service",
		Payload:     params.ServiceName,
		ServiceType: globals.LoadTypeGRPC,
		Operation:   globals.DataUpdate,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}
