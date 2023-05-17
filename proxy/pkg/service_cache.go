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

// ServiceCache 是 service 缓存的接口。
type ServiceCache interface {
	LoadService() error
	UpdateServiceCache(serviceName string, serviceType int, operation string) error
	HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error)
	GetGrpcServiceList() []*enity.ServiceDetail
	GetTcpServiceList() []*enity.ServiceDetail
}

type serviceCache struct {
	mu           sync.RWMutex
	HTTPServices *sync.Map
	TCPServices  *sync.Map
	GRPCServices *sync.Map
	sf           singleflight.Group
}

// NewServiceCache 返回一个新的 serviceCache 实例。
func NewServiceCache() *serviceCache {
	return &serviceCache{
		mu:           sync.RWMutex{},
		HTTPServices: &sync.Map{},
		TCPServices:  &sync.Map{},
		GRPCServices: &sync.Map{},
		sf:           singleflight.Group{},
	}
}

// LoadService 将所有 service 数据加载到缓存中。
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

// UpdateServiceCache 通过 serviceName ,serviceType,operation 更新 service 缓存。
func (s *serviceCache) UpdateServiceCache(serviceName string, serviceType int, operation string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 使用singleflight.Group确保同时只有一个goroutine在执行更新操作
	_, err, _ := s.sf.Do(serviceName, func() (interface{}, error) {
		tx := mysql.GetDB()

		log.Debug("debug 更新cache", zap.String("servicename", serviceName), zap.Int("serviceType", serviceType), zap.Strings("operation", strings.Split(operation, " ")))
		// 获取对应的map
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

		// 向数据库查询info
		serviceInfo, err := get(tx, &enity.ServiceInfo{ServiceName: serviceName})
		if err != nil {
			log.Debug("查询数据库失败", zap.Error(err))
			return nil, err
		}
		log.Debug("查询数据库成功", zap.Any("serviceInfo", serviceInfo))

		// 组装新的服务详情
		updatedServiceDetail, err := getServiceDetail(tx, serviceInfo)
		if err != nil {
			return nil, err
		}

		// 移除负载均衡和传输层的缓存
		LoadBalanceTransport.Remove(serviceName)
		FlowLimiter.Remove(serviceName)
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

// HTTPAccessMode 根据请求的host和path，从URL解析出服务名，通过服务名从缓存中获取对应的服务详情。
func (s *serviceCache) HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error) {
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path
	segments := strings.Split(path, "/")
	serviceName := segments[1]

	log.Debug("serviceName", zap.String("serviceName", serviceName))

	log.Debug("http host", zap.String("host", serviceName))
	s.mu.RLock()
	serviceDetail, ok := s.HTTPServices.Load(serviceName)
	s.mu.RUnlock()
	log.Debug("http serviceDetail", zap.Any("detail", serviceDetail))
	if !ok {
		return nil, fmt.Errorf("not matched service")
	}

	detail := serviceDetail.(*enity.ServiceDetail)
	if detail.HTTPRule.RuleType == globals.HTTPRuleTypeDomain {
		if detail.HTTPRule.Rule == host {
			c.Set("service", detail)
			return detail, nil
		} else {
			log.Info("host does not match domain rule", zap.String("host", host), zap.String("rule", detail.HTTPRule.Rule))
			return nil, fmt.Errorf("host does not match domain rule")
		}
	}
	if detail.HTTPRule.RuleType == globals.HTTPRuleTypePrefixURL {
		if strings.HasPrefix(path, detail.HTTPRule.Rule) {
			c.Set("service", detail)
			return detail, nil
		} else {
			log.Info("path does not match prefix URL rule", zap.String("path", path), zap.String("rule", detail.HTTPRule.Rule))
			return nil, fmt.Errorf("path does not match prefix URL rule")
		}
	}

	return nil, nil
}

// GetGrpcServiceList 遍历map获取所有的 gRPC 服务列表。
func (s *serviceCache) GetGrpcServiceList() []*enity.ServiceDetail {
	return s.getServiceListFromMap(s.GRPCServices)
}

// GetTcpServiceList 遍历map获取所有的 TCP 服务列表。
func (s *serviceCache) GetTcpServiceList() []*enity.ServiceDetail {
	return s.getServiceListFromMap(s.TCPServices)
}

// getServiceListFromMap 工具函数，工具传入的map进行遍历，返回[]*enity.ServiceDetail。
func (s *serviceCache) getServiceListFromMap(serviceMap *sync.Map) []*enity.ServiceDetail {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := []*enity.ServiceDetail{}
	serviceMap.Range(func(key, value interface{}) bool {
		serviceDetail := value.(*enity.ServiceDetail)
		list = append(list, serviceDetail)
		return true
	})
	return list

}
