package cache

import (
	"fmt"
	"gateway/enity"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"sync"

	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type AppCache interface {
	LoadAppCache() error
	UpdateAppCache(appID string) error
	GetAppList() []*enity.App
}

type appCache struct {
	AppCache     *cache.Cache
	singleFlight singleflight.Group
	rwmutex      sync.RWMutex
}

func NewAppGoCache() *appCache {
	return &appCache{
		AppCache:     cache.New(defaultExpiration, cleanupInterval),
		singleFlight: singleflight.Group{},
		rwmutex:      sync.RWMutex{},
	}
}

func (s *appCache) GetAppList() []*enity.App {
	s.rwmutex.RLock()
	items := s.AppCache.Items()
	s.rwmutex.RUnlock()

	appList := make([]*enity.App, 0)
	for _, app := range items {
		appList = append(appList, app.Object.(*enity.App))
	}

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
	a.rwmutex.Lock()
	a.AppCache.Flush()
	for _, listItem := range list {
		tmpItem := listItem
		a.AppCache.Set(tmpItem.AppID, &tmpItem, cache.DefaultExpiration)
	}
	a.rwmutex.Unlock()

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
		s.AppCache.Set(appID, appInfo, cache.DefaultExpiration)
		return nil, nil
	})
	return err
}

func (s *appCache) findAppInfoByID(appID string) (*enity.App, error) {
	s.rwmutex.RLock()
	appInfo, found := s.AppCache.Get(appID)
	s.rwmutex.RUnlock()
	if !found {
		return nil, fmt.Errorf("app not found")
	}

	return appInfo.(*enity.App), nil
}

func (s *appCache) PrintCache() {
	s.rwmutex.RLock()
	items := s.AppCache.Items()
	s.rwmutex.RUnlock()

	for key, item := range items {
		log.Debug("all service info ", zap.Any("key", key), zap.Any("item", item))
		// fmt.Printf("Key:%v, Value:%v\n", key, item)
	}
}
