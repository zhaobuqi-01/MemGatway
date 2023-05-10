package utils

import (
	"gateway/configs"
	"gateway/pkg/database/mysql"
	"gateway/pkg/database/redis"
	"gateway/pkg/log"

	"go.uber.org/zap"
)

func Init() {
	configs.Init()
	log.Init()
	mysql.Init()
	redis.Init()
}

func Cleanup() {
	Cleanuplog()
	CleanupRedis()
	CleanupMySQL()
	// flow_counter.CleanupFlowCounter()
}

func Cleanuplog() {
	if err := log.Close(); err != nil {
		log.Fatal("Failed to close log: %v", zap.Error(err))
	}
	log.Info("log closed")
}

func CleanupRedis() {
	if err := redis.CloseRedis(); err != nil {
		log.Fatal("Failed to close redis: %v", zap.Error(err))
	}
	log.Info("Redis closed")
}
func CleanupMySQL() {
	if err := mysql.CloseDB(); err != nil {
		log.Fatal("Failed to close database: %v", zap.Error(err))
	}
	log.Info("Mysql closed")
}
