package main

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/log"
	"time"
)

func main() {

	// 初始化配置
	configs.Init()
	log.Init()
	// 每 5 秒打印一次服务器配置
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		log.Debug("test logger")
		serverConf := configs.GetServerConfig()
		fmt.Printf("Server configuration: %+v\n", serverConf)
	}
	log.Info("Server started")
}
