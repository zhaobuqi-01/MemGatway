package pkg

import (
	"sync"

	"golang.org/x/time/rate"
)

type Limiter interface {
	GetLimiter(serviceName string, qps float64) (*rate.Limiter, error)
}

type flowLimiter struct {
	flowLimiterMap sync.Map
}

func NewFlowLimiter() *flowLimiter {
	return &flowLimiter{}
}

func (fl *flowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	value, ok := fl.flowLimiterMap.Load(serviceName)
	if ok {
		return value.(*rate.Limiter), nil
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))

	fl.flowLimiterMap.Store(serviceName, newLimiter)
	return newLimiter, nil
}
