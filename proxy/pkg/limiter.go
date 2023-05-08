package pkg

import (
	"gateway/proxy/limiter"
	"sync"
)

var (
	limiterOnce sync.Once
	FlowLimiter limiter.Limiter
)

func InitFlowLimiter() {
	limiterOnce.Do(
		func() {
			FlowLimiter = limiter.NewFlowLimiter()
		})
}
