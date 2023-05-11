package pkg

import (
	"sync"
)

// cache 接口组合了 AppCache 和 ServiceCache 两个接口
type cache interface {
	AppCache
	ServiceCache
}

// appCacheAndServiceCache 结构体实现了 cache 接口，包含 AppCache 和 ServiceCache
type appCacheAndServiceCache struct {
	AppCache
	ServiceCache
}

// newCache 创建一个新的 cache 实例
func newCache() cache {
	return &appCacheAndServiceCache{
		NewAppCache(),
		NewServiceCache(),
	}
}

// 定义全局变量
var (
	// Cache 提供缓存功能，包括 AppCache 和 ServiceCache
	Cache cache
	// FlowLimiter 提供限流功能
	FlowLimiter Limiter
	// LoadBalanceTransport 提供负载均衡和传输功能
	LoadBalanceTransport LoadBalanceAndTransport
	// once 用于确保全局初始化只执行一次
	once sync.Once
)

// Init 函数用于初始化全局变量，它只会被执行一次
func Init() {
	once.Do(func() {
		Cache = newCache()
		FlowLimiter = NewFlowLimiter()
		LoadBalanceTransport = NewLoadBalancerAndTransport()
	})
}
