package main

import (
	"flag"
	"gateway/backend"
	"gateway/configs"
	"gateway/http_proxy"
	"gateway/pkg/database/mysql"
	"gateway/pkg/database/redis"
	"gateway/pkg/log"
	"gateway/utils"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var run = flag.String("run", "", "input gateway or proxy")

func main() {
	flag.Parse()
	if *run == "" {
		flag.Usage()
		os.Exit(1)
	}

	gin.SetMode(gin.DebugMode)

	if *run == "gateway" {
		// 启动后台服务器
		InitAll()
		defer CleanupAll()

		go func() {
			backend.GatewayServerRun()
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		backend.GatewayServerStop()

	} else {
		// 启动代理服务器
		InitAll()
		defer CleanupAll()

		utils.ServiceManagerHandler.LoadOnce()
		utils.AppManagerHandler.LoadOnce()
		go func() {
			http_proxy.HtppProxyServerRun()
		}()
		go func() {
			http_proxy.HttpsProxyServerRun()
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		http_proxy.HttpProxyServerStop()
		http_proxy.HttpsProxyServerStop()

	}
}

func InitAll() {
	configs.Init()
	log.Init()
	mysql.Init()
	redis.Init()
}

func CleanupAll() {
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
