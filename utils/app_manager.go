package utils

import (
	"gateway/dao"
	"gateway/pkg/database/mysql"
	"net/http/httptest"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var AppManagerHandler *appManager

func init() {
	AppManagerHandler = NewAppManager()
}

type appManager struct {
	AppMap   map[string]*dao.App
	AppSlice []*dao.App
	rwmutex  sync.RWMutex
	once     sync.Once
	err      error
}

func NewAppManager() *appManager {
	return &appManager{
		AppMap:   map[string]*dao.App{},
		AppSlice: []*dao.App{},
		rwmutex:  sync.RWMutex{},
		once:     sync.Once{},
	}
}

func (s *appManager) GetAppList() []*dao.App {
	return s.AppSlice
}

func (s *appManager) LoadOnce() error {
	s.once.Do(func() {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx := mysql.GetDB()
		appQueryConditions := []func(db *gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("(name like ? or app_id like ?)", "%", "%")
			},
		}
		// 使用dao中的PageList方法获取分页的应用程序列表
		list, _, err := dao.PageList[dao.App](c, tx, appQueryConditions, 1, 99)

		if err != nil {
			s.err = err
			return
		}
		s.rwmutex.Lock()
		defer s.rwmutex.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			s.AppMap[listItem.AppID] = &tmpItem
			s.AppSlice = append(s.AppSlice, &tmpItem)
		}
	})
	return s.err
}
