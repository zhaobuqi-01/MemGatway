package logic

import (
	"gateway/internal/dao"
	"gateway/internal/dto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
func (s *serviceHttpLogic) AddHTTP(c *gin.Context, param *dto.ServiceAddHTTPInput) error {
	tx := s.db.Begin()
	if _, err := dao.Get(c, s.db, &dao.ServiceInfo{ServiceName: param.ServiceName}); err == nil {
		tx.Rollback()
		return errors.New("服务已存在")
	}

	httpUrl := &dao.HttpRule{
		RuleType: param.RuleType,
		Rule:     param.Rule,
	}

	if _, err := dao.Get(c, s.db, httpUrl); err == nil {
		tx.Rollback()
		return errors.New("服务接入前缀或域名已存在")
	}

	if len(strings.Split(param.IpList, ",")) != len(strings.Split(param.WeightList, ",")) {
		return errors.New("IP列表与权重列表数量不一致")
	}

	serviceModel := &dao.ServiceInfo{
		ServiceName: param.ServiceName,
		ServiceDesc: param.ServiceDesc,
	}

	if err := dao.Save(c, s.db, serviceModel); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "添加服务信息失败")
	}

	httpRule := &dao.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       param.RuleType,
		Rule:           param.Rule,
		NeedHttps:      param.NeedHttps,
		NeedStripUri:   param.NeedStripUri,
		NeedWebsocket:  param.NeedWebsocket,
		UrlRewrite:     param.UrlRewrite,
		HeaderTransfor: param.HeaderTransfor,
	}
	if err := dao.Save(c, s.db, httpRule); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "添加HTTP规则失败")
	}

	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          param.OpenAuth,
		BlackList:         param.BlackList,
		WhiteList:         param.WhiteList,
		ClientIPFlowLimit: param.ClientipFlowLimit,
		ServiceFlowLimit:  param.ServiceFlowLimit,
	}
	if err := dao.Update(c, s.db, accessControl); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "添加服务权限失败")
	}

	loadbalance := &dao.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              param.RoundType,
		IpList:                 param.IpList,
		WeightList:             param.WeightList,
		UpstreamConnectTimeout: param.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  param.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    param.UpstreamIdleTimeout,
		UpstreamMaxIdle:        param.UpstreamMaxIdle,
	}
	if err := dao.Update(c, s.db, loadbalance); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "添加服务负载均衡错失败")
	}
	tx.Commit()
	return nil
}

func (s *serviceHttpLogic) UpdateHTTP(c *gin.Context, params *dto.ServiceUpdateHTTPInput) error {
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return errors.New("IP列表与权重列表数量不一致")
	}
	tx := s.db.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err := dao.Get(c, s.db, serviceInfo); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "服务不存在")
	}
	serviceDetail, err := (&dao.ServiceDetail{}).ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "获取服务详情失败")
	}
	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Save(c, s.db, info); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新服务信息失败")
	}
	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err := dao.Save(c, s.db, httpRule); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新HTTP规则失败")
	}
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := dao.Save(c, s.db, accessControl); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新服务权限失败")
	}
	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = params.RoundType
	loadbalance.IpList = params.IpList
	loadbalance.WeightList = params.WeightList
	loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := dao.Save(c, s.db, loadbalance); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新服务负载均衡错失败")
	}
	tx.Commit()
	return nil
}
