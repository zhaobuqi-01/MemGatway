package database

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var redisClient *redis.Client

func Init() {

	client, err := ConnectRedis()
	if err != nil {
		logger.Fatal("Failed to connect to Redis: %v", zap.Error(err))
	}
	redisClient = client
}

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
		Addr:         fmt.Sprintf("%s", redisConfig.Addr),
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

// 获取redis连接
func GetRedisConnection() *redis.Client {
	return redisClient
}

// Close 关闭Redis连接池
func Close() error {
	err := redisClient.Close()
	if err != nil {
		return err
	}
	return nil
}
