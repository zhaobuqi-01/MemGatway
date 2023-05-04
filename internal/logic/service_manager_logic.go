package logic

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/pkg"
	"gateway/pkg/database/mysql"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceManager struct {
	ServiceMap   map[string]*dao.ServiceDetail
	ServiceSlice []*dao.ServiceDetail
	rwmutex      sync.RWMutex
	once         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*dao.ServiceDetail{},
		ServiceSlice: []*dao.ServiceDetail{},
		rwmutex:      sync.RWMutex{},
		once:         sync.Once{},
	}
}

func (s *ServiceManager) GetTcpServiceList() []*dao.ServiceDetail {
	return s.getServiceList(pkg.LoadTypeTCP)
}

func (s *ServiceManager) GetGrpcServiceList() []*dao.ServiceDetail {
	return s.getServiceList(pkg.LoadTypeGRPC)
}

func (s *ServiceManager) getServiceList(LoadType int) []*dao.ServiceDetail {
	list := []*dao.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.LoadType == LoadType {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*dao.ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//host c.Request.Host
	//path c.Request.URL.Path
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path
	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != pkg.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == pkg.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == pkg.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, fmt.Errorf("not matched service")
}

func (s *ServiceManager) LoadOnce() error {
	s.once.Do(func() {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx := mysql.GetDB()

		// 从db中分页读取基本信息
		queryConditions := []func(db *gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("(service_name like ? or service_desc like ?)", "%"+""+"%", "%"+""+"%")
			},
		}
		list, _, err := dao.PageList[dao.ServiceInfo](c, tx, queryConditions, 1, 99999)
		if err != nil {
			s.err = err
			return
		}

		s.rwmutex.Lock()
		defer s.rwmutex.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			serviceDetail, err := dao.GetServiceDetail(c, tx, &tmpItem)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})
	return s.err
}
