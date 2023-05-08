package pkg

import (
	"gateway/proxy/limiter"
	"sync"
)

var (
	limiterOnce sync.Once
	FlowLimter  limiter.Limter
)

func InitFlowLimiter() {
	limiterOnce.Do(
		func() {
			FlowLimter = limiter.NewFlowLimiter()
		})
}
