package logic

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/configs"
	"gateway/dao"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"gateway/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceInfoLogic interface {
	DeleteService(c *gin.Context, params *dto.ServiceDeleteInput) error
	GetServiceList(c *gin.Context, params *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
	GetServiceDetail(c *gin.Context, params *dto.ServiceDeleteInput) (*enity.ServiceDetail, error)
	GetServiceStat(c *gin.Context, params *dto.ServiceDeleteInput) (*dto.ServiceStatOutput, error)
}
type serviceInfoLogic struct {
	info dao.ServiceInfoService
	tcp  dao.TcpService
	grpc dao.GrpcService
	http dao.HttpService
	lb   dao.LoadBalanceService
	ac   dao.AccessControlService
	db   *gorm.DB
}

// NewserviceInfoLogic 创建serviceInfoLogic
func NewServiceInfoLogic() *serviceInfoLogic {
	return &serviceInfoLogic{
		dao.NewServiceInfoService(),
		dao.NewTcpService(),
		dao.NewGrpcService(),
		dao.NewHttpService(),
		dao.NewLoadBalanceService(),
		dao.NewAccessControlService(),
		mysql.GetDB(),
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
	list, total, err := s.info.PageList(c, s.db, queryConditions, params.PageNo, params.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get serviceInfo list")
	}

	// 格式化输出信息
	outputList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := s.info.GetServiceDetail(c, s.db, &listItem)
		if err != nil {
			log.Error("failed to get serviceDetail")
			return nil, 0, fmt.Errorf("get serviceDetail fail")
		}

		// 根据服务类型和规则生成服务地址
		serviceAddr, err := s.getServiceAddress(serviceDetail)
		if err != nil {
			return nil, 0, fmt.Errorf("get serviceAddr fail")
		}

		counter, err := globals.FlowCounter.GetCounter(serviceDetail.Info.ServiceName)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get service flow counter")
		}

		// 获取IP列表
		ipList := utils.SplitStringByComma(serviceDetail.LoadBalance.IpList)

		outputItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qps:         int(counter.QPS),
			Qpd:         int(counter.QPD),
			TotalNode:   len(ipList),
		}
		outputList = append(outputList, outputItem)
	}

	return outputList, total, nil
}

// 删除服务
func (s *serviceInfoLogic) DeleteService(c *gin.Context, params *dto.ServiceDeleteInput) error {
	// 在这里，您需要定义服务信息的实体。假设它是 `dao.ServiceInfo`
	var err error
	serviceInfo := &enity.ServiceInfo{ID: params.ID}

	serviceInfo, err = s.info.Get(c, s.db, serviceInfo)
	if err != nil {
		return fmt.Errorf("service does not exist")
	}

	// 软删除，将is_delete设置为1；如果您需要物理删除，请使用dao.Delete(c, s.db, serviceInfo)
	serviceInfo.IsDelete = 1

	err = s.info.Save(c, s.db, serviceInfo)
	if err != nil {
		return fmt.Errorf("failed to delete service")
	}

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:        "service",
		Payload:     serviceInfo.ServiceName,
		ServiceType: serviceInfo.LoadType,
		Operation:   globals.DataDelete,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
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

	serviceInfo, err = s.info.Get(c, s.db, serviceInfo)
	if err != nil {
		return nil, fmt.Errorf("service does not exist")
	}

	// 获取服务详情
	serviceDetail, err := s.info.GetServiceDetail(c, s.db, serviceInfo)
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

// 获取服务统计信息
func (s *serviceInfoLogic) GetServiceStat(c *gin.Context, params *dto.ServiceDeleteInput) (*dto.ServiceStatOutput, error) {
	serviceInfo := &enity.ServiceInfo{ID: params.ID}
	serviceDetail, err := s.info.GetServiceDetail(c, s.db, serviceInfo)
	if err != nil {
		return nil, fmt.Errorf("service does not exist")
	}
	counter, err := globals.FlowCounter.GetCounter(serviceDetail.Info.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service flow counter")
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, time.UTC)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, time.UTC)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	out := &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}
	return out, nil
}
