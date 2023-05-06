package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// ServerConfig - server configuration struct
type ServerConfig struct {
	Addr           string `mapstructure:"addr"`
	ReadTimeout    int    `mapstructure:"read_timeout"`
	WriteTimeout   int    `mapstructure:"write_timeout"`
	MaxHeaderBytes int    `mapstructure:"max_header_bytes"`
}

type GatewayServerConfig struct {
	Addr           string   `mapstructure:"addr"`
	ReadTimeout    int      `mapstructure:"read_timeout"`
	WriteTimeout   int      `mapstructure:"write_timeout"`
	MaxHeaderBytes int      `mapstructure:"max_header_bytes"`
	AllowIP        []string `mapstructure:"allow_ip"`
}

// LogConfig - 日志配置 (Log configuration)
type LogConfig struct {
	Format        string `mapstructure:"format"`
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
	SqlFile      string `mapstructure:"sql_file"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
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
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"poolsize"`
	MinIdleConns int    `mapstructure:"minidleconns"`
	DialTimeout  string `mapstructure:"dialtimeout"`
	ReadTimeout  string `mapstructure:"readtimeout"`
	WriteTimeout string `mapstructure:"writetimeout"`
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

type ClusterConfig struct {
	ClusterIp      string `mapstructure:"cluster_ip"`
	ClusterPort    string `mapstructure:"cluster_port"`
	ClusterSslPort string `mapstructure:"cluster_ssl_port"`
}

// Global configuration variables
var (
	v                   = viper.New()
	logConfig           *LogConfig
	mySQLConfig         *MySQLConfig
	redisConfig         *RedisConfig
	swaggerConfig       *SwaggerConfig
	gatewyServerConfig  *GatewayServerConfig
	httpProxyConfig     *ServerConfig
	httpsProxyConfig    *ServerConfig
	metricsServerConfig *ServerConfig
	clusterConfig       *ClusterConfig
)

func Init() {
	v.SetConfigName("config")
	v.AddConfigPath("./configs")
	v.SetConfigType("yaml") // 设置配置文件类型

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	loadConfig()
}

func loadConfig() {
	err := v.UnmarshalKey("log", &logConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'log' config: %v\n", err)
	}

	err = v.UnmarshalKey("mysql", &mySQLConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'mysql' config: %v\n", err)
	}

	err = v.UnmarshalKey("redis", &redisConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'redis' config: %v\n", err)
	}

	err = v.UnmarshalKey("swagger", &swaggerConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'swagger' config: %v\n", err)
	}

	err = v.UnmarshalKey("gateway", &gatewyServerConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'gateway' config: %v\n", err)
	}
	err = v.UnmarshalKey("http", &httpProxyConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'http' config: %v\n", err)
	}

	err = v.UnmarshalKey("https", &httpsProxyConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'https' config: %v\n", err)
	}

	err = v.UnmarshalKey("metrics", &metricsServerConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'metrics' config: %v\n", err)
	}

	err = v.UnmarshalKey("cluster", &clusterConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'cluster' config: %v\n", err)
	}
}

// 向外部暴露的函数；用于取对应的配置
func GetGatewayServerConfig() *GatewayServerConfig {
	return gatewyServerConfig
}

func GetMerricsServerConfig() *ServerConfig {
	return metricsServerConfig
}

func GetLogConfig() *LogConfig {
	return logConfig
}

func GetMysqlConfig() *MySQLConfig {
	return mySQLConfig
}

func GetRedisConfig() *RedisConfig {
	return redisConfig
}

func GetSwaggerConfig() *SwaggerConfig {
	return swaggerConfig
}

// GetHttpProxyConfig 用于获取 HTTP 代理配置
func GetHttpProxyConfig() *ServerConfig {
	return httpProxyConfig
}

// GetHttpsProxyConfig 用于获取 HTTPS 代理配置
func GetHttpsProxyConfig() *ServerConfig {
	return httpsProxyConfig
}

func GetClusterConfig() *ClusterConfig {
	return clusterConfig
}
