package database

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sync"
	"time"
)

// ConnectRedis 连接到Redis数据库 (Connect to Redis database)
func ConnectRedis() (*redis.Client, error) {
	redisConfig := configs.GetRedisConfig()
	dialTimeout, err := time.ParseDuration(redisConfig.DialTimeout)
	if err != nil {
		return nil, err
	}
	readTimeout, err := time.ParseDuration(redisConfig.ReadTimeout)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := time.ParseDuration(redisConfig.WriteTimeout)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,
		MinIdleConns: redisConfig.MinIdleConns,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

var RedisClient *redis.Client
var onceRedis sync.Once

func init() {
	onceRedis.Do(func() {
		client, err := ConnectRedis()
		if err != nil {
			logger.Fatal("Failed to connect to Redis: %v", zap.Error(err))
		}
		RedisClient = client
	})
}
