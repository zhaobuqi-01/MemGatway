package logic

import (
	"fmt"
	"gateway/configs"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceInfoLogic interface {
	ServiceDelete(c *gin.Context, param *dto.ServiceDeleteInput) error
	GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
	GetServiceDetail(c *gin.Context, param *dto.ServiceDeleteInput) (*dao.ServiceDetail, error)
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
func (s *serviceInfoLogic) GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error) {
	if s.db == nil {
		return nil, 0, fmt.Errorf("dao is not initialized")
	}

	// 从db中分页读取基本信息
	queryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
		},
	}
	list, total, err := dao.PageList[dao.ServiceInfo](c, s.db, queryConditions, param.PageNo, param.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get serviceInfo list")
	}

	// 格式化输出信息
	outputList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := dao.GetServiceDetail(c, s.db, &listItem)
		if err != nil {
			return nil, 0, fmt.Errorf("get serviceDetail fail")
		}

		// 根据服务类型和规则生成服务地址
		serviceAddr, err := s.getServiceAddress(serviceDetail)
		if err != nil {
			return nil, 0, fmt.Errorf("get serviceAddr fail")
		}

		// 获取IP列表
		ipList := pkg.SplitStringByComma(serviceDetail.LoadBalance.IpList)

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
func (s *serviceInfoLogic) ServiceDelete(c *gin.Context, param *dto.ServiceDeleteInput) error {
	// 在这里，您需要定义服务信息的实体。假设它是 `dao.ServiceInfo`
	var err error
	serviceInfo := &dao.ServiceInfo{ID: param.ID}

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

	return nil
}

// 获取服务详情
func (s *serviceInfoLogic) GetServiceDetail(c *gin.Context, param *dto.ServiceDeleteInput) (*dao.ServiceDetail, error) {
	var err error
	serviceInfo := &dao.ServiceInfo{ID: param.ID}

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
func (s *serviceInfoLogic) getServiceAddress(serviceDetail *dao.ServiceDetail) (string, error) {
	clusterIP := configs.GetStringConfig("cluster.cluster_ip")
	clusterPort := configs.GetStringConfig("cluster.cluster_port")
	clusterSSLPort := configs.GetStringConfig("cluster.cluster_ssl_port")

	switch serviceDetail.Info.LoadType {
	case pkg.LoadTypeHTTP:
		if serviceDetail.HTTPRule.RuleType == pkg.HTTPRuleTypePrefixURL {
			if serviceDetail.HTTPRule.NeedHttps == 0 {
				return fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule), nil
			}
			return fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule), nil
		}
		if serviceDetail.HTTPRule.RuleType == pkg.HTTPRuleTypeDomain {
			return serviceDetail.HTTPRule.Rule, nil
		}
		return "unknown", fmt.Errorf("unsupported load type")
	case pkg.LoadTypeTCP:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port), nil
	case pkg.LoadTypeGRPC:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port), nil
	default:
		return "unknown", fmt.Errorf("unsupported load type")
	}
}
