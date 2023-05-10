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
	UpdateServiceCache(serviceName string) error
	HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error)
	GetGrpcServiceList() []*enity.ServiceDetail
	GetTcpServiceList() []*enity.ServiceDetail
}

type serviceCache struct {
	ServiceCache *sync.Map
	sf           singleflight.Group
}

func NewServiceCache() *serviceCache {
	return &serviceCache{
		ServiceCache: &sync.Map{},
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
	list, _, err := pageList[enity.ServiceInfo](tx, serviceInfoQueryConditions, 1, 9999)
	if err != nil {
		return err
	}

	s.ServiceCache = &sync.Map{}
	for _, listItem := range list {
		tmpItem := listItem
		serviceDetail, err := getServiceDetail(tx, &tmpItem)
		if err != nil {
			return err
		}
		s.ServiceCache.Store(tmpItem.ServiceName, serviceDetail)
	}

	log.Info("load service manager successfully")
	return nil
}

func (s *serviceCache) UpdateServiceCache(serviceName string) error {
	// 使用singleflight.Group确保同时只有一个goroutine在执行更新操作
	_, err, _ := s.sf.Do(serviceName, func() (interface{}, error) {
		tx := mysql.GetDB()

		// 查询数据库获取服务详情
		serviceDetail, err := s.findServiceDetailByName(serviceName)
		if err != nil {
			return nil, err
		}

		// 提取ServiceInfo
		serviceInfo := serviceDetail.Info

		updatedServiceDetail, err := getServiceDetail(tx, serviceInfo)
		if err != nil {
			return nil, err
		}

		// 将新的服务详情设置到缓存中
		s.ServiceCache.Store(serviceName, updatedServiceDetail)
		return nil, nil
	})
	return err
}

func (s *serviceCache) HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error) {
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path

	s.ServiceCache.Range(func(key, value interface{}) bool {
		serviceDetail := value.(*enity.ServiceDetail)
		if serviceDetail.Info.LoadType != globals.LoadTypeHTTP {
			return true
		}
		if serviceDetail.HTTPRule.RuleType == globals.HTTPRuleTypeDomain {
			if serviceDetail.HTTPRule.Rule == host {
				c.Set("service", serviceDetail)
				return false
			}
		}
		if serviceDetail.HTTPRule.RuleType == globals.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceDetail.HTTPRule.Rule) {
				c.Set("service", serviceDetail)
				return false
			}
		}
		return true
	})

	serviceDetail, exists := c.Get("service")
	if !exists {
		return nil, fmt.Errorf("not matched service")
	}

	return serviceDetail.(*enity.ServiceDetail), nil
}

func (s *serviceCache) GetGrpcServiceList() []*enity.ServiceDetail {
	list := []*enity.ServiceDetail{}
	s.ServiceCache.Range(func(key, value interface{}) bool {
		serviceDetail := value.(*enity.ServiceDetail)
		if serviceDetail.Info.LoadType == globals.LoadTypeGRPC {
			list = append(list, serviceDetail)
		}
		return true
	})
	return list
}

func (s *serviceCache) GetTcpServiceList() []*enity.ServiceDetail {
	list := []*enity.ServiceDetail{}
	s.ServiceCache.Range(func(key, value interface{}) bool {
		serviceDetail := value.(*enity.ServiceDetail)
		if serviceDetail.Info.LoadType == globals.LoadTypeTCP {
			list = append(list, serviceDetail)
		}
		return true
	})
	return list
}

func (s *serviceCache) findServiceDetailByName(serviceName string) (*enity.ServiceDetail, error) {
	value, ok := s.ServiceCache.Load(serviceName)
	if !ok {
		return nil, fmt.Errorf("service not found")
	}
	return value.(*enity.ServiceDetail), nil
}

func (s *serviceCache) PrintCache() {
	s.ServiceCache.Range(func(key, value interface{}) bool {
		log.Debug("all service info", zap.Any("key", key), zap.Any("value", value))
		return true
	})
}
