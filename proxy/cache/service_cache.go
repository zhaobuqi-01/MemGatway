package cache

import (
	"fmt"
	"strings"

	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
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
	ServiceCache *cache.Cache
	sf           singleflight.Group
}

func NewServiceCache() *serviceCache {
	return &serviceCache{
		ServiceCache: cache.New(defaultExpiration, cleanupInterval),
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

	s.ServiceCache.Flush()
	for _, listItem := range list {
		tmpItem := listItem
		serviceDetail, err := getServiceDetail(tx, &tmpItem)
		if err != nil {
			return err
		}
		s.ServiceCache.Set(tmpItem.ServiceName, serviceDetail, cache.DefaultExpiration)
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
		s.ServiceCache.Set(serviceName, updatedServiceDetail, cache.DefaultExpiration)
		return nil, nil
	})
	return err
}

func (s *serviceCache) HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error) {
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path

	for key := range s.ServiceCache.Items() {
		serviceDetail, err := s.findServiceDetailByName(key)
		if err != nil {
			continue
		}

		if serviceDetail.Info.LoadType != globals.LoadTypeHTTP {
			continue
		}
		if serviceDetail.HTTPRule.RuleType == globals.HTTPRuleTypeDomain {
			if serviceDetail.HTTPRule.Rule == host {
				return serviceDetail, nil
			}
		}
		if serviceDetail.HTTPRule.RuleType == globals.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceDetail.HTTPRule.Rule) {
				return serviceDetail, nil
			}
		}
	}
	return nil, fmt.Errorf("not matched service")
}

func (s *serviceCache) GetGrpcServiceList() []*enity.ServiceDetail {
	items := s.ServiceCache.Items()

	list := []*enity.ServiceDetail{}
	for _, item := range items {
		serviceDetail := item.Object.(*enity.ServiceDetail)
		if serviceDetail.Info.LoadType == globals.LoadTypeGRPC {
			list = append(list, serviceDetail)
		}
	}
	return list
}

func (s *serviceCache) GetTcpServiceList() []*enity.ServiceDetail {
	items := s.ServiceCache.Items()

	list := []*enity.ServiceDetail{}
	for _, item := range items {
		serviceDetail := item.Object.(*enity.ServiceDetail)
		if serviceDetail.Info.LoadType == globals.LoadTypeTCP {
			list = append(list, serviceDetail)
		}
	}
	return list
}

func (s *serviceCache) findServiceDetailByName(serviceName string) (*enity.ServiceDetail, error) {
	serviceDetail, found := s.ServiceCache.Get(serviceName)
	if !found {
		return nil, fmt.Errorf("service not found")
	}
	return serviceDetail.(*enity.ServiceDetail), nil
}

func (s *serviceCache) PrintCache() {
	items := s.ServiceCache.Items()

	for key, item := range items {
		log.Debug("all service info ", zap.Any("key", key), zap.Any("item", item))
		// fmt.Printf("Key:%v, Value:%v\n", key, item)
	}
}
