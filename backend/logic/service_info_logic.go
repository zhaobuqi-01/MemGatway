package logic

import (
	"fmt"
	"gateway/backend/dao"
	"gateway/backend/dto"
	"gateway/backend/utils"
	"gateway/configs"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceInfoLogic interface {
	ServiceDelete(c *gin.Context, params *dto.ServiceDeleteInput) error
	GetServiceList(c *gin.Context, params *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
	GetServiceDetail(c *gin.Context, params *dto.ServiceDeleteInput) (*enity.ServiceDetail, error)
}
type serviceInfoLogic struct {
	db *gorm.DB
}

// NewserviceInfoLogic 创建serviceInfoLogic
func NewServiceInfoLogic(tx *gorm.DB) ServiceInfoLogic {
	return &serviceInfoLogic{
		db: tx,
	}
}

// 获取服务列表
func (s *serviceInfoLogic) GetServiceList(c *gin.Context, params *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error) {
	if s.db == nil {
		return nil, 0, fmt.Errorf("dao is not initialized")
	}

	// 从db中分页读取基本信息
	queryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
		},
	}
	list, total, err := dao.PageList[enity.ServiceInfo](c, s.db, queryConditions, params.PageNo, params.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get serviceInfo list")
	}

	// 格式化输出信息
	outputList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := dao.GetServiceDetail(c, s.db, &listItem)
		if err != nil {
			log.Error("failed to get serviceDetail")
			return nil, 0, fmt.Errorf("get serviceDetail fail")
		}

		// 根据服务类型和规则生成服务地址
		serviceAddr, err := s.getServiceAddress(serviceDetail)
		if err != nil {
			return nil, 0, fmt.Errorf("get serviceAddr fail")
		}

		// 获取IP列表
		ipList := utils.SplitStringByComma(serviceDetail.LoadBalance.IpList)

		outputItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(ipList),
		}
		outputList = append(outputList, outputItem)
	}

	return outputList, total, nil
}

// 删除服务
func (s *serviceInfoLogic) ServiceDelete(c *gin.Context, params *dto.ServiceDeleteInput) error {
	// 在这里，您需要定义服务信息的实体。假设它是 `dao.ServiceInfo`
	var err error
	serviceInfo := &enity.ServiceInfo{ID: params.ID}

	serviceInfo, err = dao.Get(c, s.db, serviceInfo)
	if err != nil {
		return fmt.Errorf("service does not exist")
	}

	// 软删除，将is_delete设置为1；如果您需要物理删除，请使用dao.Delete(c, s.db, serviceInfo)
	serviceInfo.IsDelete = 1

	err = dao.Save(c, s.db, serviceInfo)
	if err != nil {
		return fmt.Errorf("failed to delete service")
	}

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:    "service",
		Payload: serviceInfo.ServiceName,
	}
	if err := utils.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}

// 获取服务详情
func (s *serviceInfoLogic) GetServiceDetail(c *gin.Context, params *dto.ServiceDeleteInput) (*enity.ServiceDetail, error) {
	var err error
	serviceInfo := &enity.ServiceInfo{ID: params.ID}

	serviceInfo, err = dao.Get(c, s.db, serviceInfo)
	if err != nil {
		return nil, fmt.Errorf("service does not exist")
	}

	// 获取服务详情
	serviceDetail, err := dao.GetServiceDetail(c, s.db, serviceInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get serviceDetail")
	}

	return serviceDetail, nil
}

// 1、http后缀接入 clusterIP+clusterPort+path
// 2、http域名接入 domain
// 3、tcp、grpc接入 clusterIP+servicePort
// 获取服务地址
func (s *serviceInfoLogic) getServiceAddress(serviceDetail *enity.ServiceDetail) (string, error) {
	clustCfg := configs.GetClusterConfig()
	clusterIP := clustCfg.ClusterIp
	clusterPort := clustCfg.ClusterPort
	clusterSSLPort := clustCfg.ClusterSslPort

	switch serviceDetail.Info.LoadType {
	case globals.LoadTypeHTTP:
		if serviceDetail.HTTPRule.RuleType == globals.HTTPRuleTypePrefixURL {
			if serviceDetail.HTTPRule.NeedHttps == 0 {
				return fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule), nil
			}
			return fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule), nil
		}
		if serviceDetail.HTTPRule.RuleType == globals.HTTPRuleTypeDomain {
			return serviceDetail.HTTPRule.Rule, nil
		}
		return "unknown", fmt.Errorf("unsupported load type")
	case globals.LoadTypeTCP:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port), nil
	case globals.LoadTypeGRPC:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port), nil
	default:
		return "unknown", fmt.Errorf("unsupported load type")
	}
}
