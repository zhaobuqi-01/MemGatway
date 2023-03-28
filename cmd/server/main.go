package main

import (
	"gateway/configs"
	"gateway/internal/router"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gin.SetMode(configs.GetGinConfig().Mode)
	// 启动 HTTP 服务
	router.HttpServerRun()

	// 监听操作系统信号，收到信号时优雅地停止服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 停止 HTTP 服务
	router.HttpServerStop()
}
