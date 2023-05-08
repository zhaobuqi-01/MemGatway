package pkg

import (
	"gateway/proxy/cache"
	"sync"
)

type appCacheAndServiceCache struct {
	cache.AppCache
	cache.ServeiceCache
}

var (
	Cache     appCacheAndServiceCache
	cacheOnce sync.Once
)

func InitCache() {
	cacheOnce.Do(
		func() {
			Cache = appCacheAndServiceCache{
				cache.NewAppCache(),
				cache.NewServiceCache(),
			}
		})
}
