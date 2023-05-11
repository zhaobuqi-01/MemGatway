package pkg

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"

	"gateway/pkg/log"
	"sync"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

// AppCache is the interface of app cache
type AppCache interface {
	// LoadAppCache 加载所有app数据到缓存
	LoadAppCache() error
	// UpdateAppCache 通过appID,operation更新app缓存
	UpdateAppCache(appID string, operation string) error
	// GetApp 通过appID获取app
	GetApp(appID string) (*enity.App, error)
}

type appCache struct {
	AppCache     *sync.Map
	singleFlight singleflight.Group
}

// NewAppCache returns a new appCache instance
func NewAppCache() *appCache {
	return &appCache{
		AppCache:     &sync.Map{},
		singleFlight: singleflight.Group{},
	}
}

// GetApp returns app by appID
func (s *appCache) GetApp(appID string) (*enity.App, error) {
	any, ok := s.AppCache.Load(appID)
	if !ok {
		return nil, fmt.Errorf("app not found")
	}
	return any.(*enity.App), nil
}

// LoadAppCache loads all app data into cache
func (a *appCache) LoadAppCache() error {
	log.Info("start loading app to cache")
	tx := mysql.GetDB()

	appQueryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%", "%")
		},
	}
	list, err := getAll[enity.App](tx, appQueryConditions)
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

// UpdateAppCache updates app cache
func (s *appCache) UpdateAppCache(appID string, operation string) error {
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

		// 根据操作类型更新缓存
		switch operation {
		case globals.DataInsert, globals.DataUpdate:
			s.AppCache.Store(appID, appInfo)
			return nil, nil
		case globals.DataDelete:
			s.AppCache.Delete(appID)
			return nil, nil
		default:
			return nil, fmt.Errorf("invalid operation")
		}
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
