package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

type Limiter interface {
	GetLimiter(serviceName string, qps float64) (*rate.Limiter, error)
}

type FlowLimiter struct {
	flowLimiterMap sync.Map
}

type FlowLimiterItem struct {
	serviceName string
	limiter     *rate.Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{}
}

func (fl *FlowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	value, ok := fl.flowLimiterMap.Load(serviceName)
	if ok {
		return value.(*FlowLimiterItem).limiter, nil
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		serviceName: serviceName,
		limiter:     newLimiter,
	}

	fl.flowLimiterMap.Store(serviceName, item)
	return newLimiter, nil
}
