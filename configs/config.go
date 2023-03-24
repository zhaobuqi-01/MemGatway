package configs

import (
	"fmt"
	"gateway/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

// ServerConfig - server configuration struct
type ServerConfig struct {
	Addr           string   `mapstructure:"addr"`
	ReadTimeout    int      `mapstructure:"read_timeout"`
	WriteTimeout   int      `mapstructure:"write_timeout"`
	MaxHeaderBytes int      `mapstructure:"max_header_bytes"`
	AllowIP        []string `mapstructure:"allow_ip"`
}

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
	MaxIdleConns int  `mapstructure:"maxidleconns"`
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

// GetServerConfig - get server configuration
func GetServerConfig() *ServerConfig {
	v := readConfig("base", "yaml", "../../configs")

	serverConfig := &ServerConfig{}
	err := v.UnmarshalKey("server", serverConfig)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal ServerConfig: %s", err))
	}

	return serverConfig
}

// GetLogConfig - 获取日志配置 (Get logger configuration)
func GetLogConfig() *LogConfig {
	v := readConfig("logger", "yaml", "../../configs")
	logConfig := &LogConfig{}
	err := v.Unmarshal(logConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal LogConfig: %s", err))
	}
	return logConfig
}

// GetMySQLConfig - 获取MySQL配置 (Get MySQL configuration)
func GetMySQLConfig() *MySQLConfig {
	v := readConfig("mysql", "yaml", "../../configs")
	mysqlConfig := &MySQLConfig{}
	err := v.Unmarshal(mysqlConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal MySQLConfig: %s", err))
	}
	return mysqlConfig
}

// GetRedisConfig - 获取Redis配置 (Get Redis configuration)
func GetRedisConfig() *RedisConfig {
	v := readConfig("redis", "yaml", "../../configs")
	redisConfig := &RedisConfig{}
	err := v.Unmarshal(redisConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal RedisConfig: %s", err))
	}
	return redisConfig
}

// Global configuration variables
var (
	logConfig    *LogConfig
	mysqlConfig  *MySQLConfig
	redisConfig  *RedisConfig
	serverConfig *ServerConfig
)

var configOnce sync.Once

// LoadConfigurations loads configurations from the config files
// 加载配置文件中的配置
func LoadConfigurations() error {
	var err error

	configOnce.Do(func() {
		// Load server configuration
		serverConfig = GetServerConfig()

		// Load logger configurations
		logConfig = GetLogConfig()

		// Load MySQL configurations
		mysqlConfig = GetMySQLConfig()

		// Load Redis configurations
		redisConfig = GetRedisConfig()
	})

	return err
}

// 用于确保仅初始化一次
// Used to ensure that initialization occurs only once
var onceConfig sync.Once

func init() {
	// 使用 sync.Once 仅执行一次初始化
	// Use sync.Once to initialize only once
	onceConfig.Do(func() {
		err := LoadConfigurations()
		if err != nil {
			// 如果配置解析失败，则打印错误并退出
			// If the configuration parsing fails, print the error and exit
			logger.Fatal("Failed to load configurations: %v", zap.Error(err))
		}
	})
}

// 获取指定key的值
func Get(key string) interface{} {
	return viper.Get(key)
}

// 获取string类型的配置
func GetString(key string) string {
	return viper.GetString(key)
}

// 获取int类型的配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// 获取bool类型的配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// 获取string slice类型的配置
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}
