package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

const smallSliceSize = 50

// var GloablFlowLimiter *FlowLimiter

// func init() {
// 	GloablFlowLimiter = NewFlowLimiter()
// }

type Limiter interface {
	GetLimiter(serviceName string, qps float64) (*rate.Limiter, error)
}

type FlowLimiter struct {
	flowLimiterMap   sync.Map
	flowLimiterSlice []*FlowLimiterItem
	locker           sync.Mutex
}

type FlowLimiterItem struct {
	serviceName string
	limiter     *rate.Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		flowLimiterSlice: make([]*FlowLimiterItem, 0),
		locker:           sync.Mutex{},
	}
}

func (fl *FlowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	if len(fl.flowLimiterSlice) < smallSliceSize {
		for _, item := range fl.flowLimiterSlice {
			if item.serviceName == serviceName {
				return item.limiter, nil
			}
		}
	} else {
		value, ok := fl.flowLimiterMap.Load(serviceName)
		if ok {
			return value.(*FlowLimiterItem).limiter, nil
		}
	}

	fl.locker.Lock()
	defer fl.locker.Unlock()

	// Double check after acquiring the lock
	if len(fl.flowLimiterSlice) < smallSliceSize {
		for _, item := range fl.flowLimiterSlice {
			if item.serviceName == serviceName {
				return item.limiter, nil
			}
		}
	} else {
		value, ok := fl.flowLimiterMap.Load(serviceName)
		if ok {
			return value.(*FlowLimiterItem).limiter, nil
		}
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		serviceName: serviceName,
		limiter:     newLimiter,
	}

	// map存储全部的限流器
	fl.flowLimiterMap.Store(serviceName, item)
	// slice存储最近50个服务的限流器
	if len(fl.flowLimiterSlice) < smallSliceSize {
		fl.flowLimiterSlice = append(fl.flowLimiterSlice, item)
	}

	return newLimiter, nil
}
