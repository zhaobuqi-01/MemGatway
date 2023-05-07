package counter

import (
	"sync/atomic"
	"time"
)

type Counter struct {
	value int64
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Add(num int64) {
	atomic.AddInt64(&c.value, num)
}

func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *Counter) Reset(t time.Duration) {
	ticker := time.NewTicker(24 * t)
	defer ticker.Stop()

	for range ticker.C {
		atomic.StoreInt64(&c.value, 0)
	}
}
