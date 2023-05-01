package main

import (
	"fmt"
	"gateway/configs"
	"time"
)

func main() {
	// 初始化配置
	configs.Init()

	// 每 5 秒打印一次服务器配置
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		serverConf := configs.GetServerConfig()
		fmt.Printf("Server configuration: %+v\n", serverConf)
	}
}
