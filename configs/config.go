package configs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadConfigurations loads configurations from the config files
// init 初始化配置
func Init() {
	// Set configuration file path
	setConfigPath()

	// Load server configuration
	getCfg()
}

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
	Level         string `mapstructure:"level"`
	Filename      string `mapstructure:"filename"`
	ErrorFilename string `mapstructure:"error_filename"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}

// MySQLConfig - MySQL配置 (MySQL configuration)
type MySQLConfig struct {
	SqlFile      string
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	Collation    string `mapstructure:"collation"`
	ParseTime    bool   `mapstructure:"parsetime"`
	MaxIdleConns int    `mapstructure:"maxidleconns"`
	MaxOpenConns int    `mapstructure:"maxopenconns"`
}

// RedisConfig - Redis配置 (Redis configuration)
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int    `mapstructure:"minidleconns"`
	DialTimeout  string `mapstructure:"dialtimeout"`
	ReadTimeout  string `mapstructure:"readtimeout"`
	WriteTimeout string `mapstructure:"writetimeout"`
}

type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

// SwaggerConfig - Swagger配置 (Swagger configuration)
type SwaggerConfig struct {
	Version     string   `mapstructure:"version"`
	Host        string   `mapstructure:"host"`
	BasePath    string   `mapstructure:"base_path"`
	Schemes     []string `mapstructure:"schemes"`
	Title       string   `mapstructure:"title"`
	Description string   `mapstructure:"description"`
}

var ConfigPath string // 硬编码的配置文件路径

func setConfigPath() {
	if os.Getenv("GATEWAY_CONFIG_PATH") != "" {
		// WSL2
		ConfigPath = os.Getenv("GATEWAY_CONFIG_PATH")
	} else {
		// Windows
		ConfigPath = "E:\\project\\Api-Gateway\\configs"
	}
}

// readConfig - 读取配置文件 (Read configuration file)
func readConfig() *viper.Viper {
	// 实例化viper
	v := viper.New()
	// 配置文件名称（无扩展名）
	v.SetConfigName("config")
	// 配置文件类型，如果配置文件的名称
	v.SetConfigType("yaml")
	// 查找配置文件所在的路径
	v.AddConfigPath(ConfigPath)

	// 查找并读取配置文件
	err := v.ReadInConfig()

	// 处理读取配置文件的错误
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}

	return v
}

func getConfig[T GinConfig | LogConfig | RedisConfig | MySQLConfig | ServerConfig | SwaggerConfig](describe string, configType *T) *T {
	v := readConfig()
	err := v.UnmarshalKey(describe, configType)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal RedisConfig: %s", err))
	}
	return configType
}

// Global configuration variables
var (
	logConfig     *LogConfig
	mysqlConfig   *MySQLConfig
	redisConfig   *RedisConfig
	serverConfig  *ServerConfig
	ginConfig     *GinConfig
	swaggerConfig *SwaggerConfig
)

// 向外部暴露的函数；用于取对应的配置
func GetServerConfig() *ServerConfig {
	return serverConfig
}

func GetLogConfig() *LogConfig {
	return logConfig
}

func GetMysqlConfig() *MySQLConfig {
	return mysqlConfig
}

func GetRedisConfig() *RedisConfig {
	return redisConfig
}

func GetGinConfig() *GinConfig {
	return ginConfig
}

// 向外部暴露的函数；用于取对应的配置
func GetSwaggerConfig() *SwaggerConfig {
	return swaggerConfig
}

func getCfg() {
	// Load server configuration
	serverConfig = getConfig("server", new(ServerConfig))

	// Load logger configurations
	logConfig = getConfig("log", new(LogConfig))

	// Load MySQL configurations
	mysqlConfig = getConfig("mysql", new(MySQLConfig))

	// Load Redis configurations
	redisConfig = getConfig("redis", new(RedisConfig))

	ginConfig = getConfig("gin", new(GinConfig))

	swaggerConfig = getConfig("swagger", new(SwaggerConfig))
}
