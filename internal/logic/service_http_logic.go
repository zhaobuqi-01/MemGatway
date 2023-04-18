package logic

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type serviceHttpLogic struct {
	db *gorm.DB
}

// NewserviceHttpLogic 创建serviceHttpLogic
func NewServiceHttpLogic(tx *gorm.DB) *serviceHttpLogic {
	return &serviceHttpLogic{
		db: tx,
	}
}

// 添加HTTP服务
func (s *serviceHttpLogic) AddHTTP(c *gin.Context, params *dto.ServiceAddHTTPInput) error {
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return fmt.Errorf("IP列表与权重列表数量不一致")
	}

	tx := s.db.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err := dao.Get(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		return fmt.Errorf("服务已存在")
	}

	httpUrl := &dao.HttpRule{
		RuleType: params.RuleType,
		Rule:     params.Rule,
	}

	if _, err := dao.Get(c, tx, httpUrl); err == nil {
		tx.Rollback()
		return fmt.Errorf("服务接入前缀或域名已存在")
	}
	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}

	if err := dao.Save(c, tx, serviceModel); err != nil {
		tx.Rollback()
		return fmt.Errorf("添加服务信息失败")
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
		return fmt.Errorf("添加HTTP规则失败")
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
		return fmt.Errorf("添加服务权限失败")
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
		return fmt.Errorf("添加服务负载均衡错失败")
	}
	tx.Commit()
	return nil
}

func (s *serviceHttpLogic) UpdateHTTP(c *gin.Context, paramss *dto.ServiceUpdateHTTPInput) error {
	if len(strings.Split(paramss.IpList, ",")) != len(strings.Split(paramss.WeightList, ",")) {
		return fmt.Errorf("IP列表与权重列表数量不一致")
	}

	tx := s.db.Begin()

	serviceInfo, err := dao.Get(c, tx, &dao.ServiceInfo{ServiceName: paramss.ServiceName})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("服务不存在")
	}

	serviceDetail, err := (&dao.ServiceDetail{}).ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("服务不存在")
	}

	info := serviceDetail.Info
	info.ServiceDesc = paramss.ServiceDesc
	if err := dao.Update(c, tx, info); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新服务描述失败")
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = paramss.NeedHttps
	httpRule.NeedStripUri = paramss.NeedStripUri
	httpRule.NeedWebsocket = paramss.NeedWebsocket
	httpRule.UrlRewrite = paramss.UrlRewrite
	httpRule.HeaderTransfor = paramss.HeaderTransfor
	if err := dao.Update(c, tx, httpRule); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新HTTP规则失败")
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = paramss.OpenAuth
	accessControl.BlackList = paramss.BlackList
	accessControl.WhiteList = paramss.WhiteList
	accessControl.ClientIPFlowLimit = paramss.ClientipFlowLimit
	accessControl.ServiceFlowLimit = paramss.ServiceFlowLimit
	if err := dao.Update(c, tx, accessControl); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新服务权限失败")
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
		return fmt.Errorf("更新服务负载均衡错失败")
	}

	tx.Commit()

	return nil
}
