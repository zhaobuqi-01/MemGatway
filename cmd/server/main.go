package main

import (
	"gateway/configs"
	"gateway/internal/router"
	"gateway/pkg/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitAll()
	gin.SetMode(configs.GetGinConfig().Mode)
	// 启动 HTTP 服务
	router.HttpServerRun()

	// 监听操作系统信号，收到信号时优雅地停止服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 停止 HTTP 服务
	router.HttpServerStop()

	// 等待HttpServerStop执行完毕，执行清理操作
	utils.CleanupMySQL()
	utils.CleanupRedis()
	utils.CleanupLogger()
}
