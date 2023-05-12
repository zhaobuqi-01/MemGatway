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

// AppCache 是 app 缓存的接口。
type AppCache interface {
	// LoadAppCache 将所有app数据加载到缓存中。
	LoadAppCache() error
	// UpdateAppCache 通过appID和operation更新app缓存。
	UpdateAppCache(appID string, operation string) error
	// GetApp 通过appID获取app。
	GetApp(appID string) (*enity.App, error)
}

// appCache 结构体实现了 AppCache 接口。
type appCache struct {
	AppCache     *sync.Map
	singleFlight singleflight.Group
}

// NewAppCache 返回一个新的 appCache 实例。
func NewAppCache() *appCache {
	return &appCache{
		AppCache:     &sync.Map{},
		singleFlight: singleflight.Group{},
	}
}

// GetApp 通过 appID 返回 app。
func (s *appCache) GetApp(appID string) (*enity.App, error) {
	any, ok := s.AppCache.Load(appID)
	if !ok {
		return nil, fmt.Errorf("app not found")
	}
	return any.(*enity.App), nil
}

// LoadAppCache 将所有 app 数据加载到缓存中。
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

	// 将新数据加载到缓存中
	a.AppCache = &sync.Map{}
	for _, listItem := range list {
		tmpItem := listItem
		a.AppCache.Store(tmpItem.AppID, &tmpItem)
	}

	log.Info("load app to cache successfully")
	return nil
}

// UpdateAppCache 通过 appID 和 operation 更新 app 缓存。
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

// findAppInfoByID 通过 appID 查找 app 信息。
func (s *appCache) findAppInfoByID(appID string) (*enity.App, error) {
	appInfo, found := s.AppCache.Load(appID)
	if !found {
		return nil, fmt.Errorf("app not found")
	}

	return appInfo.(*enity.App), nil
}
