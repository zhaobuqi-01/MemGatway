package flow_counter

import (
	"sync"
	"time"
)

type FlowCounter interface {
	GetCounter(counterName string) (*redisFlowCounter, error)
}

type flowCounter struct {
	redisFlowCountMap sync.Map
}

func NewFlowCounter() *flowCounter {
	return &flowCounter{
		redisFlowCountMap: sync.Map{},
	}
}

func (counter *flowCounter) GetCounter(counterName string) (*redisFlowCounter, error) {
	value, ok := counter.redisFlowCountMap.Load(counterName)
	if ok {
		return value.(*redisFlowCounter), nil
	}

	newCounter := NewRedisFlowCountService(counterName, 1*time.Second)
	counter.redisFlowCountMap.Store(counterName, newCounter)
	return newCounter, nil
}
