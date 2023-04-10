package service

import (
	"fmt"
	"gateway/configs"
	"gateway/internal/dto"
	"gateway/internal/entity"
	"gateway/internal/pkg"
	"gateway/internal/repository"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceInfoService interface {
	GetServiceAddress(serviceDetail *entity.ServiceDetail) (string, error)
	GetIPList(c *gin.Context, data *entity.LoadBalance) []string
	GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error)
}

type serviceInfoService struct {
	repo repository.ServiceInfo
}

func NewServiceInfoService(db *gorm.DB) *serviceInfoService {
	var repo repository.ServiceInfo

	if db != nil {
		repo = repository.NewServiceInfoRepo(db)
	}

	return &serviceInfoService{
		repo: repo,
	}
}

func (s *serviceInfoService) GetServiceList(c *gin.Context, param *dto.ServiceListInput) ([]dto.ServiceListItemOutput, int64, error) {
	if s.repo == nil {
		return nil, 0, errors.New("repository is not initialized")
	}

	// 从db中分页读取基本信息
	list, total, err := s.repo.PageList(c, param)
	if err != nil {
		return nil, 0, err
	}

	// 格式化输出信息
	outputList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := s.repo.ServiceDetail(c, &listItem)
		if err != nil {
			return nil, 0, err
		}

		// 根据服务类型和规则生成服务地址
		serviceAddr, err := s.GetServiceAddress(serviceDetail)
		if err != nil {
			return nil, 0, err
		}

		// 获取IP列表
		ipList := s.GetIPList(c, serviceDetail.LoadBalance)

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
func (s *serviceInfoService) GetServiceAddress(serviceDetail *entity.ServiceDetail) (string, error) {
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

func (s *serviceInfoService) GetIPList(c *gin.Context, data *entity.LoadBalance) []string {
	return strings.Split(data.IpList, ",")
}

func (s *serviceInfoService) Delete(c *gin.Context, param *dto.ServiceDeleteInput) error {
	// 在这里，您需要定义服务信息的实体。假设它是 `entity.ServiceInfo`
	var err error
	serviceInfo := &entity.ServiceInfo{ID: param.ID}

	serviceInfo, err = s.repo.Get(c, serviceInfo)
	if err != nil {
		return err
	}

	serviceInfo.IsDelete = 1

	err = s.repo.Update(c, serviceInfo)
	if err != nil {
		return err
	}

	return nil
}
