package logic

import (
	"fmt"

	"gateway/backend/dao"
	"gateway/backend/dto"
	"gateway/backend/utils"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"

	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TcpServiceLogic interface {
	AddTCP(c *gin.Context, param *dto.ServiceAddTcpInput) error
	UpdateTCP(c *gin.Context, param *dto.ServiceUpdateTcpInput) error
}

type tcpServiceLogic struct {
	db *gorm.DB
}

// NewTcpServiceLogic 创建tcpServiceLogic
func NewTcpServiceLogic(tx *gorm.DB) TcpServiceLogic {
	return &tcpServiceLogic{
		db: tx,
	}
}

// AddTCP 添加TCP服务
func (s *tcpServiceLogic) AddTCP(c *gin.Context, params *dto.ServiceAddTcpInput) error {
	// 检查服务名是否被占用
	infoSearch := &enity.ServiceInfo{ServiceName: params.ServiceName, IsDelete: 0}
	if info, err := dao.Get(c, s.db, infoSearch); err != gorm.ErrRecordNotFound {
		if err == nil && info != nil {
			return fmt.Errorf("the TCP service name already exists, please change the service name")
		}
		return fmt.Errorf("an error occurred while querying the TCP service name")
	}

	// 检查端口是否被占用
	if _, err := dao.Get(c, s.db, &enity.TcpRule{Port: params.Port}); err == nil {
		return fmt.Errorf("the port already exists, please change the port")
	}
	if _, err := dao.Get(c, s.db, &enity.GrpcRule{Port: params.Port}); err == nil {
		return fmt.Errorf("the port already exists, please change the port")
	}

	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()
	info := &enity.ServiceInfo{
		LoadType:    globals.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add TCP service information")
	}
	loadBalance := &enity.LoadBalance{
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
	tcpRule := &enity.TcpRule{
		ServiceID: info.ID,
		Port:      params.Port,
	}
	if err := dao.Save(c, tx, tcpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add TCP service rule information")
	}
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
		return fmt.Errorf("failed to add TCP service permission information")
	}
	tx.Commit()

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:    "service",
		Payload: params.ServiceName,
	}
	if err := utils.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}

// UpdateTCP 更新TCP服务
func (s *tcpServiceLogic) UpdateTCP(c *gin.Context, params *dto.ServiceUpdateTcpInput) error {
	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()

	detail, err := dao.GetServiceDetail(c, s.db, &enity.ServiceInfo{ID: params.ID})
	if err != nil {
		return fmt.Errorf("TCP service does not exist")
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save TCP service description")
	}

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
		return fmt.Errorf("failed to Save TCP service load balancing information")
	}

	tcpRule := &enity.TcpRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := dao.Save(c, tx, tcpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save TCP service rule information")
	}

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
		return fmt.Errorf("failed to Save TCP service permission information")
	}

	tx.Commit()

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:    "service",
		Payload: params.ServiceName,
	}
	if err := utils.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}
