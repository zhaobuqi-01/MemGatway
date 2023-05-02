package configs

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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

var (
	v           = viper.New()
	mutex       sync.RWMutex
	reloadMutex sync.Mutex
	reloadTimer *time.Timer
	reloadDelay = 5 * time.Second // 设置防抖动延迟时间
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

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		scheduleReloadConfig()
	})
}

func loadConfig() {
	mutex.Lock()
	defer mutex.Unlock()

	err := v.UnmarshalKey("server", &serverConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'server' config: %v\n", err)
	}

	err = v.UnmarshalKey("log", &logConfig)
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

	err = v.UnmarshalKey("gin", &ginConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'gin' config: %v\n", err)
	}

	err = v.UnmarshalKey("swagger", &swaggerConfig)
	if err != nil {
		log.Printf("Error unmarshalling 'swagger' config: %v\n", err)
	}

	log.Println("Reloaded configuration")
}

// 使用回调函数可以解决循环引用的问题，因为回调函数是在配置文件发生变化时被触发的，
// 而不是在包的初始化过程中被调用的。这意味着在回调函数中调用其他模块的函数不会导致循环引用的问题。
// 当我们在一个模块中调用另一个模块的函数时，如果两个模块都引用了对方，就会形成循环引用。
// 这会导致编译错误或者运行时错误，因为两个模块都无法正确初始化。
// 但是，如果我们使用回调函数，那么这种情况就不会发生，
// 因为回调函数是在配置文件发生变化时被触发的。这意味着每个模块都可以独立地初始化，
// 而不需要在初始化过程中相互依赖。当配置文件发生变化时，每个模块都会被重新初始化，
// 并被重新注册为回调函数，以确保它们能够正确地响应配置变化。
var reloadCallbacks []func()

// RegisterReloadCallback 注册配置重新加载回调函数
func RegisterReloadCallback(callback func()) {
	reloadCallbacks = append(reloadCallbacks, callback)
}

func scheduleReloadConfig() {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	if reloadTimer != nil {
		reloadTimer.Stop()
	}

	reloadTimer = time.AfterFunc(reloadDelay, func() {
		log.Println("Reloading configuration after debounce delay")
		loadConfig()
		// 调用所有已注册的回调函数
		for _, callback := range reloadCallbacks {
			callback()
		}
	})
}

// Global configuration variables
var (
	logConfig     *LogConfig
	mySQLConfig   *MySQLConfig
	redisConfig   *RedisConfig
	serverConfig  *ServerConfig
	ginConfig     *GinConfig
	swaggerConfig *SwaggerConfig
)

// 向外部暴露的函数；用于取对应的配置
func GetServerConfig() *ServerConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return serverConfig
}

func GetLogConfig() *LogConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return logConfig
}

func GetMysqlConfig() *MySQLConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return mySQLConfig
}

func GetRedisConfig() *RedisConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return redisConfig
}

func GetGinConfig() *GinConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return ginConfig
}

// 向外部暴露的函数；用于取对应的配置
func GetSwaggerConfig() *SwaggerConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return swaggerConfig
}

func GetStringConfig(key string) string {
	mutex.RLock()
	defer mutex.RUnlock()
	return v.GetString(key)
}

func GetIntConfig(key string) int {
	mutex.RLock()
	defer mutex.RUnlock()
	return v.GetInt(key)
}

func GetBoolConfig(key string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return v.GetBool(key)
}

func GetSliceConfig(key string) []string {
	mutex.RLock()
	defer mutex.RUnlock()
	return v.GetStringSlice(key)
}
