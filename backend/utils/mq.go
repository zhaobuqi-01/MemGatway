package utils

import (
	"gateway/mq"
	"gateway/pkg/database/redis"
	"sync"
)

var (
	MessageQueue *mq.MessageQueue
	once         sync.Once
)

func InitMq() {
	once.Do(func() {
		MessageQueue = mq.Default(redis.GetRedisConnection())
	})
}
