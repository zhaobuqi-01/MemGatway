package service

import (
	"errors"
	"fmt"
	"gateway/configs"
	"gateway/internal/model"
	"gateway/internal/pkg"
	"gateway/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServiceInfoService interface {
	GetServiceAddress(serviceDetail *model.ServiceDetail) (string, error)
	GetIPList(c *gin.Context, data *model.LoadBalance) []string
}

type serviceInfoService struct {
	repo repository.ServiceInfo
}

func NewServiceInfoService(repo repository.ServiceInfo) ServiceInfoService {
	return &serviceInfoService{
		repo: repo,
	}
}

func (s *serviceInfoService) GetServiceAddress(serviceDetail *model.ServiceDetail) (string, error) {
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
		return serviceDetail.HTTPRule.Rule, nil
	case pkg.LoadTypeTCP:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port), nil
	case pkg.LoadTypeGRPC:
		return fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port), nil
	default:
		return "unknown", errors.New("unsupported load type")
	}
}

func (s *serviceInfoService) GetIPList(c *gin.Context, data *model.LoadBalance) []string {
	return strings.Split(data.IpList, ",")
}
