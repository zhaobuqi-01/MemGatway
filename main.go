package main

import (
	"flag"
	"gateway/backend"
	"gateway/http_proxy"
	"gateway/metrics"
	"gateway/utils"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var run = flag.String("run", "", "input gateway or proxy")

func main() {
	flag.Parse()
	if *run == "" {
		flag.Usage()
		os.Exit(1)
	}
	metrics.RecordSystemMetrics()
	gin.SetMode(gin.DebugMode)

	if *run == "gateway" {
		// 启动后台服务器
		utils.InitAll()
		defer utils.CleanupAll()

		go func() {
			backend.GatewayServerRun()
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		backend.GatewayServerStop()

	} else {
		// 启动代理服务器
		utils.InitAll()
		defer utils.CleanupAll()

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
