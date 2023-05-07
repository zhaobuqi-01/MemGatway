package logic

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/dao"
	"strings"

	"github.com/gin-gonic/gin"
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
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err := dao.Get(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		return fmt.Errorf("HTTP service already exists")
	}

	httpUrl := &dao.HttpRule{
		RuleType: params.RuleType,
		Rule:     params.Rule,
	}

	if _, err := dao.Get(c, tx, httpUrl); err == nil {
		tx.Rollback()
		return fmt.Errorf("HTTP service access prefix or domain name already exists")
	}
	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := dao.Save(c, tx, serviceModel); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add HTTP service information")
	}

	httpRule := &dao.HttpRule{
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

	accessControl := &dao.AccessControl{
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

	loadbalance := &dao.LoadBalance{
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
	return nil
}

func (s *httpServiceLogic) UpdateHTTP(c *gin.Context, paramss *dto.ServiceUpdateHTTPInput) error {
	if len(strings.Split(paramss.IpList, ",")) != len(strings.Split(paramss.WeightList, ",")) {
		return fmt.Errorf("the IP list is inconsistent with the number of weight lists")
	}

	tx := s.db.Begin()

	serviceInfo, err := dao.Get(c, tx, &dao.ServiceInfo{ServiceName: paramss.ServiceName})
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
	info.ServiceDesc = paramss.ServiceDesc
	if err := dao.Update(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update HTTP service description")
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = paramss.NeedHttps
	httpRule.NeedStripUri = paramss.NeedStripUri
	httpRule.NeedWebsocket = paramss.NeedWebsocket
	httpRule.UrlRewrite = paramss.UrlRewrite
	httpRule.HeaderTransfor = paramss.HeaderTransfor
	if err := dao.Update(c, tx, httpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update HTTP service rules")
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = paramss.OpenAuth
	accessControl.BlackList = paramss.BlackList
	accessControl.WhiteList = paramss.WhiteList
	accessControl.ClientIPFlowLimit = paramss.ClientipFlowLimit
	accessControl.ServiceFlowLimit = paramss.ServiceFlowLimit
	if err := dao.Update(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update HTTP service permissions")
	}

	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = paramss.RoundType
	loadbalance.IpList = paramss.IpList
	loadbalance.WeightList = paramss.WeightList
	loadbalance.UpstreamConnectTimeout = paramss.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = paramss.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = paramss.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = paramss.UpstreamMaxIdle
	if err := dao.Update(c, tx, loadbalance); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update HTTP service load balancing error")
	}

	tx.Commit()

	return nil
}
