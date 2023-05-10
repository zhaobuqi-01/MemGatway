```go
package pkg

import (
	"fmt"

	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"

	"net/http/httptest"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ServiceManagerHandler *serviceManager

var serviceManagerOnce sync.Once

func NewServiceManagerHandlerOnce() {
	serviceManagerOnce.Do(func() {
		ServiceManagerHandler = NewServiceManager()
	})
}

type serviceManager struct {
	ServiceMap   map[string]*enity.ServiceDetail
	ServiceSlice []*enity.ServiceDetail
	rwmutex      sync.RWMutex
}

func NewServiceManager() *serviceManager {
	return &serviceManager{
		ServiceMap:   map[string]*enity.ServiceDetail{},
		ServiceSlice: []*enity.ServiceDetail{},
		rwmutex:      sync.RWMutex{},
	}
}

func (s *serviceManager) GetTcpServiceList() []*enity.ServiceDetail {
	list := []*enity.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.LoadType == globals.LoadTypeTCP {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *serviceManager) GetGrpcServiceList() []*enity.ServiceDetail {
	list := []*enity.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.LoadType == globals.LoadTypeGRPC {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *serviceManager) HTTPAccessMode(c *gin.Context) (*enity.ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//host c.Request.Host
	//path c.Request.URL.Path
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path
	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != globals.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == globals.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == globals.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, fmt.Errorf("not matched service")
}

func (s *serviceManager) Load(logMsg string) error {
	log.Info("start loading service manager", zap.String("logMsg", logMsg))
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	tx := mysql.GetDB()

	// 从db中分页读取基本信息
	serviceInfoQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(service_name like ? or service_desc like ?)", "%", "%")
		},
	}
	list, _, err := pageList[enity.ServiceInfo](c, tx, serviceInfoQueryConditions, 1, 9999)
	if err != nil {
		return err
	}

	newServiceMap := make(map[string]*enity.ServiceDetail)
	newServiceSlice := []*enity.ServiceDetail{}

	for _, listItem := range list {
		tmpItem := listItem
		serviceDetail, err := getServiceDetail(c, tx, &tmpItem)
		if err != nil {
			return err
		}
		newServiceMap[listItem.ServiceName] = serviceDetail
		newServiceSlice = append(newServiceSlice, serviceDetail)
	}

	s.rwmutex.Lock()
	s.ServiceMap = newServiceMap
	s.ServiceSlice = newServiceSlice
	s.rwmutex.Unlock()

	log.Info("load service manager successfully", zap.String("logMsg", logMsg))
	return nil
}

package pkg

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"net/http/httptest"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	AppManagerHandler *appManager
	appManagerOnce    sync.Once
)

func NewAppManagerHandlerOnce() {
	appManagerOnce.Do(func() {
		AppManagerHandler = NewAppManager()
	})
}

type appManager struct {
	AppMap   map[string]*enity.App
	AppSlice []*enity.App
	rwmutex  sync.RWMutex
}

func NewAppManager() *appManager {
	return &appManager{
		AppMap:   map[string]*enity.App{},
		AppSlice: []*enity.App{},
		rwmutex:  sync.RWMutex{},
	}
}

func (s *appManager) GetAppList() []*enity.App {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()
	return s.AppSlice
}

func (s *appManager) Load(logMsg string) error {
	log.Info(fmt.Sprintf("start  %s loading app list", logMsg))
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	tx := mysql.GetDB()
	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}
	// 使用dao中的PageList方法获取分页的应用程序列表
	list, _, err := pageList[enity.App](c, tx, appQueryConditions, 1, 99)
	if err != nil {
		return err
	}

	// Create temporary data structures for new data
	newAppMap := map[string]*enity.App{}
	newAppSlice := []*enity.App{}

	for _, listItem := range list {
		tmpItem := listItem
		newAppMap[listItem.AppID] = &tmpItem
		newAppSlice = append(newAppSlice, &tmpItem)
	}

	s.rwmutex.Lock()
	s.AppMap = newAppMap
	s.AppSlice = newAppSlice
	s.rwmutex.Unlock()

	log.Info("load app list successfully")
	return nil
}

```

我的tcp反向代理，grpc反向代理，http反向代理都会使用缓存在内存的消息；tcp，grpc，http，https服务器同时运行在一台机器上，怎么进行缓存，可以使得性能最优