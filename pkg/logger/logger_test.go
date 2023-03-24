// logger/logger_test.go

package logger_test

import (
	"gateway/configs"
	"gateway/pkg/logger"
	"testing"
)

func TestInitLogger(t *testing.T) {
	// Load configurations first
	err := configs.LoadConfigurations()
	if err != nil {
		t.Fatalf("Failed to load configurations: %v", err)
	}

	err = logger.InitLogger()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Test logger with different log levels
	logger.Info("This is an info log")
	logger.Warn("This is a warning log")
	logger.Error("This is an error log")
}
