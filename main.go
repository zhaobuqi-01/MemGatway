package main

import (
	"fmt"
	"gateway/configs"
)

func main() {
	configs.Init()
	fmt.Println(configs.GetStringConfig("cluster.cluster_ip"),
		configs.GetStringConfig("cluster.cluster_port"),
		configs.GetStringConfig("cluster.cluster_ssl_port"))
}
