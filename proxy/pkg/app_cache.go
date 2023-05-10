package pkg

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/database/mysql"

	"gateway/pkg/log"
	"sync"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type AppCache interface {
	LoadAppCache() error
	UpdateAppCache(appID string) error
	GetAppList() []*enity.App
}

type appCache struct {
	AppCache     *sync.Map
	singleFlight singleflight.Group
}

func NewAppCache() *appCache {
	return &appCache{
		AppCache:     &sync.Map{},
		singleFlight: singleflight.Group{},
	}
}

func (s *appCache) GetAppList() []*enity.App {
	appList := make([]*enity.App, 0)
	s.AppCache.Range(func(_, value interface{}) bool {
		appList = append(appList, value.(*enity.App))
		return true
	})
	return appList
}

func (a *appCache) LoadAppCache() error {
	log.Info("start loading app to cache")
	tx := mysql.GetDB()

	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}
	// 使用dao中的PageList方法获取分页的应用程序列表
	list, _, err := pageList[enity.App](tx, appQueryConditions, 1, 9999)
	if err != nil {
		return err
	}

	// Load new data into cache
	a.AppCache = &sync.Map{}
	for _, listItem := range list {
		tmpItem := listItem
		a.AppCache.Store(tmpItem.AppID, &tmpItem)
	}

	log.Info("load app to cache successfully")
	return nil
}

func (s *appCache) UpdateAppCache(appID string) error {
	_, err, _ := s.singleFlight.Do(appID, func() (interface{}, error) {
		tx := mysql.GetDB()

		// 查询数据库获取应用程序详情
		appInfo, err := s.findAppInfoByID(appID)
		if err != nil {
			return nil, err
		}

		// 查询数据库获得app
		appInfo, err = get(tx, appInfo)
		if err != nil {
			return nil, err
		}

		// 将新的app设置到缓存中
		s.AppCache.Store(appID, appInfo)
		return nil, nil
	})
	return err
}

func (s *appCache) findAppInfoByID(appID string) (*enity.App, error) {
	appInfo, found := s.AppCache.Load(appID)
	if !found {
		return nil, fmt.Errorf("app not found")
	}

	return appInfo.(*enity.App), nil
}
