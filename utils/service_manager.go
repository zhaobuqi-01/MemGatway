package utils

import (
	"fmt"
	"gateway/dao"
	"gateway/pkg/database/mysql"

	"net/http/httptest"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ServiceManagerHandler *serviceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

type serviceManager struct {
	ServiceMap   map[string]*dao.ServiceDetail
	ServiceSlice []*dao.ServiceDetail
	rwmutex      sync.RWMutex
	once         sync.Once
	err          error
}

func NewServiceManager() *serviceManager {
	return &serviceManager{
		ServiceMap:   map[string]*dao.ServiceDetail{},
		ServiceSlice: []*dao.ServiceDetail{},
		rwmutex:      sync.RWMutex{},
		once:         sync.Once{},
	}
}

func (s *serviceManager) GetTcpServiceList() []*dao.ServiceDetail {
	list := []*dao.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.LoadType == LoadTypeTCP {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *serviceManager) GetGrpcServiceList() []*dao.ServiceDetail {
	list := []*dao.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.LoadType == LoadTypeGRPC {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *serviceManager) HTTPAccessMode(c *gin.Context) (*dao.ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//host c.Request.Host
	//path c.Request.URL.Path
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path
	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, fmt.Errorf("not matched service")
}

func (s *serviceManager) LoadOnce() error {
	s.once.Do(func() {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx := mysql.GetDB()

		// 从db中分页读取基本信息
		serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("(service_name like ? or service_desc like ?)", "%", "%")
			},
		}
		list, _, err := dao.PageList[dao.ServiceInfo](c, tx, serviceInfoQueryConditions, 1, 9999)
		if err != nil {
			s.err = err
			return
		}

		s.rwmutex.Lock()
		defer s.rwmutex.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			// log.Debug("service info", zap.Any("service info", tmpItem))
			serviceDetail, err := dao.GetServiceDetail(c, tx, &tmpItem)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})

	// log.Debug("ServiceSlice", zap.Any("ServiceSlice", s.ServiceSlice))
	return s.err
}
