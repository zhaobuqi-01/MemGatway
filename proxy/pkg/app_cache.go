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
	mu           sync.RWMutex
	AppCache     *sync.Map
	singleFlight singleflight.Group
}

// NewAppCache 返回一个新的 appCache 实例。
func NewAppCache() *appCache {
	return &appCache{
		mu:           sync.RWMutex{},
		AppCache:     &sync.Map{},
		singleFlight: singleflight.Group{},
	}
}

// GetApp 通过 appID 返回 app。
func (s *appCache) GetApp(appID string) (*enity.App, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()
	tx := mysql.GetDB()

	// 查询数据库获得app
	appInfo, err := get(tx, &enity.App{AppID: appID})
	if err != nil {
		return err
	}

	// 根据操作类型更新缓存
	switch operation {
	case globals.DataInsert, globals.DataUpdate:
		s.AppCache.Store(appID, appInfo)
		return nil
	case globals.DataDelete:
		s.AppCache.Delete(appID)
		return nil
	default:
		return fmt.Errorf("invalid operation")
	}
}

// findAppInfoByID 通过 appID 查找 app 信息。
func (s *appCache) findAppInfoByID(appID string) (*enity.App, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	appInfo, found := s.AppCache.Load(appID)
	if !found {
		return nil, fmt.Errorf("app not found")
	}

	return appInfo.(*enity.App), nil
}
