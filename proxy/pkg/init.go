package pkg

import (
	"sync"
)

type cache interface {
	AppCache
	ServiceCache
}

type appCacheAndServiceCache struct {
	AppCache
	ServiceCache
}

func NewCache() cache {
	return &appCacheAndServiceCache{
		NewAppCache(),
		NewServiceCache(),
	}
}

var (
	Cache                cache
	FlowLimiter          Limiter
	LoadBalanceTransport LoadBalanceAndTransport
	once                 sync.Once
)

func Init() {
	once.Do(func() {
		Cache = NewCache()
		FlowLimiter = NewFlowLimiter()
		LoadBalanceTransport = NewLoadBalancerAndTransport()
	})
}
