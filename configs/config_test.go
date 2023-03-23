package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySQLConfigParsing(t *testing.T) {
	mysqlConfig := GetMySQLConfig()

	assert.NotEmpty(t, mysqlConfig.Host, "MySQL Host should not be empty")
	assert.NotEmpty(t, mysqlConfig.User, "MySQL User should not be empty")
	assert.NotEmpty(t, mysqlConfig.Password, "MySQL Password should not be empty")
	assert.NotEmpty(t, mysqlConfig.DBName, "MySQL DBName should not be empty")
	assert.NotEmpty(t, mysqlConfig.Charset, "MySQL Charset should not be empty")
	assert.NotZero(t, mysqlConfig.Port, "MySQL Port should not be zero")
	assert.NotZero(t, mysqlConfig.MaxIdleConns, "MySQL MaxIdleConns should not be zero")
	assert.NotZero(t, mysqlConfig.MaxOpenConns, "MySQL MaxOpenConns should not be zero")
}

func TestRedisConfigParsing(t *testing.T) {
	redisConfig := GetRedisConfig()

	assert.NotEmpty(t, redisConfig.Host, "Redis Host should not be empty")
	assert.NotEmpty(t, redisConfig.DialTimeout, "Redis DialTimeout should not be empty")
	assert.NotEmpty(t, redisConfig.ReadTimeout, "Redis ReadTimeout should not be empty")
	assert.NotEmpty(t, redisConfig.WriteTimeout, "Redis WriteTimeout should not be empty")
	assert.NotZero(t, redisConfig.Port, "Redis Port should not be zero")
	//assert.NotZero(t, redisConfig.DB, "Redis DB should not be zero")
	assert.NotZero(t, redisConfig.PoolSize, "Redis PoolSize should not be zero")
	assert.NotZero(t, redisConfig.MinIdleConns, "Redis MinIdleConns should not be zero")
}

func TestLogConfigParsing(t *testing.T) {
	logConfig := GetLogConfig()

	assert.NotEmpty(t, logConfig.Filename, "Log Filename should not be empty")
	assert.NotZero(t, logConfig.MaxSize, "Log MaxSize should not be zero")
	assert.NotZero(t, logConfig.MaxAge, "Log MaxAge should not be zero")
	assert.NotZero(t, logConfig.MaxBackups, "Log MaxBackups should not be zero")
	assert.NotEmpty(t, logConfig.Level, "Log Level should not be empty")
}
