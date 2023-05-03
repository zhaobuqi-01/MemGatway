package main

import (
	"flag"
	"gateway/internal/logic"
	"gateway/internal/pkg"
	"gateway/internal/server"
	"os"
	"os/signal"
	"syscall"
)

var endpoint = flag.String("endpoint", "", "input endpoint backend or server")

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "backend" {
		// 启动后台服务器
		pkg.InitAll()
		defer pkg.CleanupAll()

		server.HtppServerRun()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		server.HttpServerStop()

	} else {
		// 启动代理服务器
		pkg.InitAll()
		defer pkg.CleanupAll()

		logic.ServiceManagerHandler.LoadOnce()
		logic.AppManagerHandler.LoadOnce()

		server.HtppProxyServerRun()
		server.HttpsProxyServerRun()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		server.HttpProxyServerStop()
		server.HttpsProxyServerStop()

	}
}
