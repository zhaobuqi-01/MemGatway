package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

// LogConfig - 日志配置 (Log configuration)
type LogConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Level      string
}

// MySQLConfig - MySQL配置 (MySQL configuration)
type MySQLConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string `mapstructure:"dbname"`
	Charset      string
	ParseTime    bool `mapstructure:"parsetime"`
	MaxIdleConns int  `mapstructure:"max_idleconns"`
	MaxOpenConns int  `mapstructure:"maxopenconns"`
}

// RedisConfig - Redis配置 (Redis configuration)
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int    `mapstructure:"minidleconns"`
	DialTimeout  string `mapstructure:"dialtimeout"`
	ReadTimeout  string `mapstructure:"readtimeout"`
	WriteTimeout string `mapstructure:"writetimeout"`
}

// readConfig - 读取配置文件 (Read configuration file)
func readConfig(configName string, configType string, configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading %s config file: %s", configName, err))
	}
	return v
}

// GetLogConfig - 获取日志配置 (Get log configuration)
func GetLogConfig() *LogConfig {
	v := readConfig("log", "yaml", "./configs")
	logConfig := &LogConfig{}
	err := v.Unmarshal(logConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal LogConfig: %s", err))
	}
	return logConfig
}

// GetMySQLConfig - 获取MySQL配置 (Get MySQL configuration)
func GetMySQLConfig() *MySQLConfig {
	v := readConfig("mysql", "yaml", "./configs")
	mysqlConfig := &MySQLConfig{}
	err := v.Unmarshal(mysqlConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal MySQLConfig: %s", err))
	}
	return mysqlConfig
}

// GetRedisConfig - 获取Redis配置 (Get Redis configuration)
func GetRedisConfig() *RedisConfig {
	v := readConfig("redis", "yaml", "./configs")
	redisConfig := &RedisConfig{}
	err := v.Unmarshal(redisConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal RedisConfig: %s", err))
	}
	return redisConfig
}
