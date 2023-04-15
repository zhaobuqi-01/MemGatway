package logic

import (
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
			return errors.New("服务名已存在,请更换服务名")
		}
		return errors.Wrap(err, "查询服务名时发生错误")
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := dao.Get(c, s.db, tcpRuleSearch); err == nil {
		return errors.Wrap(err, "端口已存在,请更换端口")
	}

	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := dao.Get(c, s.db, grpcRuleSearch); err == nil {
		return errors.Wrap(err, "端口已存在,请更换端口")
	}

	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return errors.New("ip列表与权重列表数量不一致")
	}

	tx := s.db.Begin()

	info := &dao.ServiceInfo{
		LoadType:    pkg.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := dao.Save(c, tx, info); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "添加服务信息失败")
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
		return errors.Wrap(err, "添加服务负载均衡失败")
	}

	grpcRule := &dao.GrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := dao.Save(c, tx, grpcRule); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "添加服务grpc规则失败")
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
		return errors.Wrap(err, "添加服务权限失败")
	}

	tx.Commit()

	return nil
}

func (s *servcieGrpcLogic) UpdateGrpc(c *gin.Context, params *dto.ServiceUpdateGrpcInput) error {
	// ip列表与权重列表数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		return errors.New("ip列表与权重列表数量不一致")
	}

	tx := s.db.Begin()

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := (&dao.ServiceDetail{}).ServiceDetail(c, s.db, service)
	if err != nil {
		return errors.Wrap(err, "服务不存在")
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := dao.Update(c, tx, info); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新服务信息失败")
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
		return errors.Wrap(err, "更新服务负载均衡失败")
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
		return errors.Wrap(err, "更新服务grpc规则失败")
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
		return errors.Wrap(err, "更新服务权限失败")
	}

	tx.Commit()

	return nil
}
