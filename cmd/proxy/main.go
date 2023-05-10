package main

import (
	"encoding/json"
	"gateway/globals"
	"gateway/mq"
	"gateway/pkg/database/redis"
	"gateway/pkg/log"
	httpRouter "gateway/proxy/http_proxy/router"
	"gateway/proxy/pkg"
	"os"
	"os/signal"
	"syscall"

	Init "gateway/init"

	"go.uber.org/zap"
)

func main() {
	pkg.Init()
	globals.Init()
	Init.Init()
	defer Init.Cleanup()

	// Load data from the database
	if err := pkg.Cache.LoadService(); err != nil {
		log.Fatal("failed to load service manager", zap.Error(err))
	}
	if err := pkg.Cache.LoadAppCache(); err != nil {
		log.Fatal("failed to load app manager", zap.Error(err))
	}

	// Create a message queue instance
	messageQueue := mq.Default(redis.GetRedisConnection())
	// Subscribe to data change channel and reload data
	err := messageQueue.Subscribe(globals.DataChange, true, func(channel string, message []byte) {
		// parse the message
		var dataChangeMsg globals.DataChangeMessage
		err := json.Unmarshal(message, &dataChangeMsg)
		if err != nil {
			log.Error("failed to unmarshal message", zap.Error(err))
			return
		}
		// handle the message based on the type
		switch dataChangeMsg.Type {
		case "app":
			appID := dataChangeMsg.Payload
			//  update app cache
			if err := pkg.Cache.UpdateAppCache(appID); err != nil {
				log.Error("failed to update app cache", zap.Error(err))
				return
			}
		case "service":
			serviceName := dataChangeMsg.Payload
			// update service cache
			if err := pkg.Cache.UpdateServiceCache(serviceName); err != nil {
				log.Error("failed to update service cache", zap.Error(err))
				return
			}
		default:
			log.Warn("unknown message type", zap.String("type", dataChangeMsg.Type))
		}
		log.Info("subscribed to data change messages", zap.String("channel", channel), zap.String("message", string(message)))
	})
	if err != nil {
		log.Fatal("failed to subscribe to data change messages", zap.Error(err))
	}

	go func() {
		httpRouter.HtppProxyServerRun()
	}()

	go func() {
		httpRouter.HttpsProxyServerRun()
	}()

	// // 启动GrpcProxyServer
	// go func() {
	// 	// grpc_proxy.Run()
	// 	grpcRouter.GrpcProxyServerRun()
	// }()

	// // run tcp proxy server
	// go func() {
	// 	// tcp_proxy.Run()
	// 	tcpRouter.TcpProxyServerRun()
	// }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// stop http proxy server
	httpRouter.HttpProxyServerStop()
	// // stop grpc proxy server
	// grpcRouter.GrpcProxyServerStop()
	// // stop tcp proxy server
	// tcpRouter.TcpProxyServerStop()
}
