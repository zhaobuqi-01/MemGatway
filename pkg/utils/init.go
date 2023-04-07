package utils

import (
	"gateway/configs"
	mysql "gateway/pkg/database/mysql"
	redis "gateway/pkg/database/redis"
	"gateway/pkg/logger"

	"go.uber.org/zap"
)

func InitAll() {
	configs.Init()
	logger.Init()
	mysql.InitDB()
	redis.Init()
}

func CleanupLogger() {
	if err := logger.Close(); err != nil {
		logger.Fatal("Failed to close logger: %v", zap.Error(err))
	}
	logger.Info("Logger closed")
}

func CleanupRedis() {
	if err := redis.Close(); err != nil {
		logger.Fatal("Failed to close redis: %v", zap.Error(err))
	}
	logger.Info("Redis closed")
}
func CleanupMySQL() {
	if err := mysql.CloseDB(); err != nil {
		logger.Fatal("Failed to close database: %v", zap.Error(err))
	}
	logger.Info("Database closed")
}
