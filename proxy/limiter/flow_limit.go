package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

const smallSliceSize = 50

// var GloablFlowLimiter *flowLimiter

// func init() {
// 	GloablFlowLimiter = NewFlowLimiter()
// }

type Limiter interface {
	GetLimiter(serviceName string, qps float64) (*rate.Limiter, error)
}

type flowLimiter struct {
	flowLimiterMap   sync.Map
	flowLimiterSlice []*flowLimiterItem
	locker           sync.Mutex
}

type flowLimiterItem struct {
	serviceName string
	limiter     *rate.Limiter
}

func NewFlowLimiter() *flowLimiter {
	return &flowLimiter{
		flowLimiterSlice: make([]*flowLimiterItem, 0),
		locker:           sync.Mutex{},
	}
}

func (fl *flowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	if len(fl.flowLimiterSlice) < smallSliceSize {
		for _, item := range fl.flowLimiterSlice {
			if item.serviceName == serviceName {
				return item.limiter, nil
			}
		}
	} else {
		value, ok := fl.flowLimiterMap.Load(serviceName)
		if ok {
			return value.(*flowLimiterItem).limiter, nil
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
			return value.(*flowLimiterItem).limiter, nil
		}
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &flowLimiterItem{
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
