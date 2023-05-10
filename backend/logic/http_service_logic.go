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

type HttpServiceLogic interface {
	AddHTTP(c *gin.Context, param *dto.ServiceAddHTTPInput) error
	UpdateHTTP(c *gin.Context, param *dto.ServiceUpdateHTTPInput) error
}

type httpServiceLogic struct {
	db *gorm.DB
}

// NewHttpServiceLogic 创建serviceHttpLogic
func NewHttpServiceLogic(tx *gorm.DB) *httpServiceLogic {
	return &httpServiceLogic{
		db: tx,
	}
}

// 添加HTTP服务
func (s *httpServiceLogic) AddHTTP(c *gin.Context, params *dto.ServiceAddHTTPInput) error {
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()
	serviceInfo := &enity.ServiceInfo{ServiceName: params.ServiceName}
	if _, err := dao.Get(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		return fmt.Errorf("HTTP service already exists")
	}

	httpUrl := &enity.HttpRule{
		RuleType: params.RuleType,
		Rule:     params.Rule,
	}

	if _, err := dao.Get(c, tx, httpUrl); err == nil {
		tx.Rollback()
		return fmt.Errorf("HTTP service access prefix or domain name already exists")
	}
	serviceModel := &enity.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := dao.Save(c, tx, serviceModel); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add HTTP service information")
	}

	httpRule := &enity.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := dao.Save(c, tx, httpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add HTTP service information")
	}

	accessControl := &enity.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := dao.Save(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add HTTP service permission")
	}

	loadbalance := &enity.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := dao.Save(c, tx, loadbalance); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add HTTP service load balancing error")
	}
	tx.Commit()

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:    "service",
		Payload: params.ServiceName,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

func (s *httpServiceLogic) UpdateHTTP(c *gin.Context, params *dto.ServiceUpdateHTTPInput) error {
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()

	serviceInfo, err := dao.Get(c, tx, &enity.ServiceInfo{ServiceName: params.ServiceName})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HTTP service does not exist")
	}

	serviceDetail, err := dao.GetServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HTTP service does not exist")
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save HTTP service description")
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err := dao.Save(c, tx, httpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save HTTP service rules")
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := dao.Save(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save HTTP service permissions")
	}

	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = params.RoundType
	loadbalance.IpList = params.IpList
	loadbalance.WeightList = params.WeightList
	loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := dao.Save(c, tx, loadbalance); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Save HTTP service load balancing error")
	}

	tx.Commit()

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:    "service",
		Payload: params.ServiceName,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}
