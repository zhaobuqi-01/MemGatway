package init

import (
	"gateway/pkg/database"
	log "gateway/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sync"
)

var RedisClient *redis.Client
var onceRedis sync.Once

func init() {
	onceRedis.Do(func() {
		client, err := database.ConnectRedis()
		if err != nil {
			log.Fatal("Failed to connect to Redis: %v", zap.Error(err))
		}
		RedisClient = client
	})
}
