package pkg

import (
	"gateway/configs"

	"gateway/pkg/database"
	"gateway/pkg/logger"

	"go.uber.org/zap"
)

func InitAll() {
	configs.Init()
	logger.Init()
	database.InitDB()
	database.InitRedis()
}

func CleanupLogger() {
	if err := logger.Close(); err != nil {
		logger.Fatal("Failed to close logger: %v", zap.Error(err))
	}
	logger.Info("Logger closed")
}

func CleanupRedis() {
	if err := database.CloseRedis(); err != nil {
		logger.Fatal("Failed to close redis: %v", zap.Error(err))
	}
	logger.Info("Redis closed")
}
func CleanupMySQL() {
	if err := database.CloseDB(); err != nil {
		logger.Fatal("Failed to close database: %v", zap.Error(err))
	}
	logger.Info("Mysql closed")
}
