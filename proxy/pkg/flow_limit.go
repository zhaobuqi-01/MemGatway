package pkg

import (
	"sync"

	"golang.org/x/time/rate"
)

// Limiter 接口定义了获取限流器的方法
type Limiter interface {
	// GetLimiter 根据服务名和每秒请求数(QPS)获取对应的限流器
	GetLimiter(serviceName string, qps float64) (*rate.Limiter, error)
}

// flowLimiter 结构体实现了Limiter接口，并包含了一个sync.Map存储各服务的限流器
type flowLimiter struct {
	flowLimiterMap sync.Map
}

// NewFlowLimiter 创建并返回一个新的flowLimiter实例
func NewFlowLimiter() *flowLimiter {
	return &flowLimiter{}
}

// GetLimiter 实现了Limiter接口中的GetLimiter方法，该方法首先检查映射中是否已经有对应服务的限流器，
// 如果有则返回，如果没有则新建一个限流器并存入sync.Map中
func (fl *flowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	value, ok := fl.flowLimiterMap.Load(serviceName)
	if ok {
		return value.(*rate.Limiter), nil
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))

	fl.flowLimiterMap.Store(serviceName, newLimiter)
	return newLimiter, nil
}
