package main

import (
	"context"
	"encoding/json"
	"gateway/configs"
	"gateway/globals"
	"gateway/mq"
	"gateway/pkg/database/redis"
	"gateway/pkg/log"
	"gateway/proxy/http_proxy/router"
	"gateway/proxy/pkg"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	Init "gateway/init"

	"go.uber.org/zap"
)

var (
	htppsProxySrv *http.Server
	htppProxySrv  *http.Server
)

func htppProxyServerRun() {
	// 初始化路由
	r := router.InitRouter()

	serverConfig := configs.GetHttpProxyConfig()
	log.Info("httpServerConfig", zap.Any("serverConfig", serverConfig))

	htppProxySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}

	log.Info("HtppProxyServer start running", zap.String("addr", serverConfig.Addr))
	if err := htppProxySrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: ", zap.String("httpProyxAddr", serverConfig.Addr), zap.Error(err))
	}
	log.Info("HtppProxyServer is running", zap.String("addr", serverConfig.Addr))

}
func httpProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := htppProxySrv.Shutdown(ctx); err != nil {
		log.Fatal("HtppProxyServer forced to shutdown:  ", zap.Error(err))
	}
	log.Info("HtppProxyServer exiting")
}

func httpsProxyServerRun() {
	// 初始化路由
	r := router.InitRouter()

	serverConfig := configs.GetHttpsProxyConfig()
	log.Info("httpsServerConfig", zap.Any("serverConfig", serverConfig))

	htppsProxySrv = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}

	log.Info("HtppsProxyServer start running", zap.String("addr", serverConfig.Addr))
	if err := htppsProxySrv.ListenAndServeTLS("cmd/http_proxy/cert_file/server.crt", "cmd/http_proxy/cert_file/server.key"); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: ", zap.String("httpsProxyAddr", serverConfig.Addr), zap.Error(err))
	}
	log.Info("HtppsProxyServer is running", zap.String("addr", serverConfig.Addr))

}

func httpsProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := htppsProxySrv.Shutdown(ctx); err != nil {
		log.Fatal("HtppsProxyServer forced to shutdown:  ", zap.Error(err))
	}
	log.Info("HtppsProxyServer exiting")
}

func main() {
	Init.InitAll()
	defer Init.CleanupAll()

	pkg.InitBalanceAndTransport()
	pkg.InitCache()
	pkg.InitFlowLimiter()
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
		htppProxyServerRun()
	}()

	go func() {
		httpsProxyServerRun()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	httpProxyServerStop()
	httpsProxyServerStop()
}
