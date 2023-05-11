package pkg

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type ServiceCache interface {
	LoadService() error
	UpdateServiceCache(serviceName string, serviceType int, operation string) error
	HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error)
	GetGrpcServiceList() []*enity.ServiceDetail
	GetTcpServiceList() []*enity.ServiceDetail
}

type serviceCache struct {
	HTTPServices *sync.Map
	TCPServices  *sync.Map
	GRPCServices *sync.Map
	sf           singleflight.Group
}

func NewServiceCache() *serviceCache {
	return &serviceCache{
		HTTPServices: &sync.Map{},
		TCPServices:  &sync.Map{},
		GRPCServices: &sync.Map{},
		sf:           singleflight.Group{},
	}
}

func (s *serviceCache) LoadService() error {
	log.Info("start loading service manager")
	tx := mysql.GetDB()

	// 从db中分页读取基本信息
	serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%", "%")
		},
	}
	list, err := getAll[enity.ServiceInfo](tx, serviceInfoQueryConditions)
	if err != nil {
		return err
	}

	for _, listItem := range list {
		tmpItem := listItem
		serviceDetail, err := getServiceDetail(tx, &tmpItem)
		if err != nil {
			return err
		}

		switch serviceDetail.Info.LoadType {
		case globals.LoadTypeHTTP:
			log.Debug("http serviceName", zap.String("service", tmpItem.ServiceName))
			s.HTTPServices.Store(tmpItem.ServiceName, serviceDetail)
		case globals.LoadTypeTCP:
			s.TCPServices.Store(tmpItem.ServiceName, serviceDetail)
		case globals.LoadTypeGRPC:
			s.GRPCServices.Store(tmpItem.ServiceName, serviceDetail)
		}
	}

	log.Info("load service manager successfully")
	return nil
}

func (s *serviceCache) UpdateServiceCache(serviceName string, serviceType int, operation string) error {
	// 使用singleflight.Group确保同时只有一个goroutine在执行更新操作
	_, err, _ := s.sf.Do(serviceName, func() (interface{}, error) {
		tx := mysql.GetDB()

		// 查询数据库获取服务详情
		var serviceMap *sync.Map
		switch serviceType {
		case globals.LoadTypeHTTP:
			serviceMap = s.HTTPServices
		case globals.LoadTypeTCP:
			serviceMap = s.TCPServices
		case globals.LoadTypeGRPC:
			serviceMap = s.GRPCServices
		default:
			return nil, fmt.Errorf("invalid service type")
		}

		serviceDetail, err := s.findServiceDetailByName(serviceName, serviceMap)
		if err != nil {
			return nil, err
		}

		// 提取ServiceInfo
		serviceInfo := serviceDetail.Info

		updatedServiceDetail, err := getServiceDetail(tx, serviceInfo)
		if err != nil {
			return nil, err
		}

		// 将新的服务详情设置到缓存
		switch operation {
		case globals.DataInsert, globals.DataUpdate:
			serviceMap.Store(serviceName, updatedServiceDetail)
			return nil, nil
		case globals.DataDelete:
			serviceMap.Delete(serviceName)
			return nil, nil
		default:
			return nil, fmt.Errorf("invalid operation")
		}
	})
	return err
}

func (s *serviceCache) HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error) {
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path
	segments := strings.Split(path, "/")
	serviceName := segments[1]

	log.Debug("http host", zap.String("host", serviceName))
	serviceDetail, ok := s.HTTPServices.Load(serviceName)
	log.Debug("http serviceDetail", zap.Any("detail", serviceDetail))
	if !ok {
		return nil, fmt.Errorf("not matched service")
	}

	detail := serviceDetail.(*enity.ServiceDetail)
	if detail.HTTPRule.RuleType == globals.HTTPRuleTypeDomain {
		if detail.HTTPRule.Rule == host {
			c.Set("service", detail)
			return detail, nil
		}
	}
	if detail.HTTPRule.RuleType == globals.HTTPRuleTypePrefixURL {
		if strings.HasPrefix(path, detail.HTTPRule.Rule) {
			c.Set("service", detail)
			return detail, nil
		}
	}

	return nil, fmt.Errorf("not matched service")
}

func (s *serviceCache) GetGrpcServiceList() []*enity.ServiceDetail {
	return s.getServiceListFromMap(s.GRPCServices)
}

func (s *serviceCache) GetTcpServiceList() []*enity.ServiceDetail {
	return s.getServiceListFromMap(s.TCPServices)
}

func (s *serviceCache) getServiceListFromMap(serviceMap *sync.Map) []*enity.ServiceDetail {
	list := []*enity.ServiceDetail{}
	serviceMap.Range(func(key, value interface{}) bool {
		serviceDetail := value.(*enity.ServiceDetail)
		list = append(list, serviceDetail)
		return true
	})
	return list

}

func (s *serviceCache) findServiceDetailByName(serviceName string, serviceMap *sync.Map) (*enity.ServiceDetail, error) {
	value, ok := serviceMap.Load(serviceName)
	if !ok {
		return nil, fmt.Errorf("service not found")
	}
	return value.(*enity.ServiceDetail), nil
}
