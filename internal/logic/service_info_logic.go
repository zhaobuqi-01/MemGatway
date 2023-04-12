package logic

import (
	"fmt"
	"gateway/configs"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type serviceInfoLogic struct {
	db *gorm.DB
}

// NewserviceInfoLogic 创建serviceInfoLogic
func NewServiceInfoLogic(tx *gorm.DB) *serviceInfoLogic {
	return &serviceInfoLogic{
		db: tx,
	}
}

// 获取服务列表
func (s *serviceInfoLogic) GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error) {
	if s.db == nil {
		return nil, 0, errors.New("dao is not initialized")
	}

	// 从db中分页读取基本信息
	list, total, err := dao.PageList[dao.ServiceInfo](c, s.db, param.Info, param.PageNo, param.PageSize)
	if err != nil {
		return nil, 0, errors.Wrap(err, "PageList(c, param.Info, param.PageNo, param.PageSize)")
	}

	// 格式化输出信息
	outputList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := (&dao.ServiceDetail{}).ServiceDetail(c, s.db, &listItem)
		if err != nil {
			return nil, 0, errors.Wrap(err, "s.getServiceDetail(c, &listItem)")
		}

		// 根据服务类型和规则生成服务地址
		serviceAddr, err := s.getServiceAddress(serviceDetail)
		if err != nil {
			return nil, 0, errors.Wrap(err, "s.getServiceAddress(serviceDetail)")
		}

		// 获取IP列表
		ipList := s.getIPList(c, serviceDetail.LoadBalance)

		outputItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
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

// 在这里放置重构后的代码
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
		return "unknown", errors.New("unsupported load type")
	case pkg.LoadTypeTCP:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port), nil
	case pkg.LoadTypeGRPC:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port), nil
	default:
		return "unknown", errors.New("unsupported load type")
	}
}

// 获取IP列表
func (s *serviceInfoLogic) getIPList(c *gin.Context, data *dao.LoadBalance) []string {
	return strings.Split(data.IpList, ",")
}

// 删除服务
func (s *serviceInfoLogic) Delete(c *gin.Context, param *dto.ServiceDeleteInput) error {
	// 在这里，您需要定义服务信息的实体。假设它是 `dao.ServiceInfo`
	var err error
	serviceInfo := &dao.ServiceInfo{ID: param.ID}

	serviceInfo, err = dao.Get(c, s.db, serviceInfo)
	if err != nil {
		return errors.Wrap(err, "Get(c, serviceInfo)")
	}

	serviceInfo.IsDelete = 1

	err = dao.Update(c, s.db, serviceInfo)
	if err != nil {
		return errors.Wrap(err, "Update(c, serviceInfo)")
	}

	return nil
}
