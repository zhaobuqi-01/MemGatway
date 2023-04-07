package router

import (
	"context"
	"gateway/configs"
	"gateway/pkg/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	r := InitRouter()
	serverConfig := configs.GetServerConfig()

	HttpSrvHandler = &http.Server{
		Addr:           serverConfig.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(serverConfig.MaxHeaderBytes),
	}
	go func() {
		logger.Info("HttpServerRun", zap.String("addr", configs.GetServerConfig().Addr))
		if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HttpServerRun", zap.String("addr", configs.GetServerConfig().Addr), zap.Error(err))
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		logger.Error("HttpServerStop  ", zap.Error(err))
	}
	logger.Info("HttpServerStop stopped")
}
