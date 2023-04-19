package redis

import (
	"context"
	"gateway/configs"
	"gateway/pkg/log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
	redisConfig *configs.RedisConfig
	ctx         = context.Background()
)

type (
	Pipeliner     = redis.Pipeliner
	Client        = redis.Client
	Cmd           = redis.Cmd
	Cmdable       = redis.Cmdable
	RedisPipeline = redis.Pipeline
	StringCmd     = redis.StringCmd
	IntCmd        = redis.IntCmd
	DurationCmd   = redis.DurationCmd
)

// InitRedis 初始化Redis数据库 (Initialize Redis database)
func InitRedis() {
	redisConfig = configs.GetRedisConfig()
	client, err := connectRedis()
	if err != nil {
		log.Fatal("Failed to connect to Redis: %v", zap.Error(err))
	}
	redisClient = client
}

// ConnectRedis 连接到Redis数据库 (Connect to Redis database)
func connectRedis() (*redis.Client, error) {
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
		Addr:         redisConfig.Addr,
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

// GetRedisConnection 获取redis客户端
func GetRedisConnection() *redis.Client {
	return redisClient
}

// Close 关闭redis客户端
func CloseRedis() error {
	err := redisClient.Close()
	if err != nil {
		return err
	}
	return nil
}

// Set 设置键值
func Set(key string, value interface{}, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

// SetNX 设置键值，如果键不存在
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return redisClient.SetNX(ctx, key, value, expiration).Result()
}

// GetKey 获取键值
func Get(key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

// Delete 删除键
func Delete(key string) (int64, error) {
	return redisClient.Del(ctx, key).Result()
}

// Expire 设置键的过期时间
func Expire(key string, expiration time.Duration) (bool, error) {
	return redisClient.Expire(ctx, key, expiration).Result()
}

// Incr 增加键的值
func Incr(key string) error {
	return redisClient.Incr(ctx, key).Err()
}

// GetInt64 获取 int64 类型的值
func GetInt64(key string) (int64, error) {
	val, err := Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// IncrWithExpire 对指定的key执行自增操作，并设置过期时间
// key: 需要自增的键
// expiration: 过期时间
//
// Example:
//
//	count, err := IncrWithExpire("counter", 10*time.Second)
//	if err != nil {
//		fmt.Println("IncrWithExpire error:", err)
//	}
func IncrWithExpire(key string, expiration time.Duration) (int64, error) {
	pipe := redisClient.Pipeline()
	incr := pipe.Incr(ctx, key)
	expire := pipe.Expire(ctx, key, expiration)

	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}

	return incr.Val(), expire.Err()
}

// Pipeline 执行管道操作
// pip: 一个或多个要在管道中执行的函数
//
// Example:
//
//	err := Pipeline(
//		func(pipe redis.Pipeliner) error {
//			return pipe.Incr(ctx, "counter").Err()
//		},
//		func(pipe redis.Pipeliner) error {
//			return pipe.Expire(ctx, "counter", 10*time.Second).Err()
//		},
//	)
//	if err != nil {
//		fmt.Println("Pipeline error:", err)
//	}
func Pipeline(pip ...func(pipe redis.Pipeliner) error) error {
	pipe := redisClient.Pipeline()
	for _, f := range pip {
		f(pipe)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// Do 执行redis命令
// commandName: Redis命令的名称，如 "SET", "GET" 等
// args: Redis命令的参数
// 返回执行结果和可能出现的错误
//
// Example:
//
//	result, err := Do("GET", "key")
//	if err != nil {
//		fmt.Println("Do error:", err)
//	}
func Do(commandName string, args ...interface{}) (interface{}, error) {
	return redisClient.Do(ctx, append([]interface{}{commandName}, args...)...).Result()
}
