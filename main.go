package main

import (
	"fmt"
	"gateway/configs"
)

func main() {
	fmt.Printf("test %s %s", configs.GetRedisConfig().Addr, configs.GetServerConfig().Addr)
	addr := configs.GetRedisConfig().Addr
	fmt.Print(addr)
}
