package router

import (
	"context"
	"fmt"
	"gateway/configs"
	"gateway/pkg/logger"
	"net/http"
	"time"
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
		logger.Default().Info(fmt.Sprintf("HttpServerRun: %s", configs.GetServerConfig().Addr))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			logger.Default().Error(fmt.Sprintf("HttpServerRun: %s err: %v", configs.GetServerConfig().Addr, err))
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		logger.Default().Error(fmt.Sprintf("HttpServerStop err: %v", err))
	}
	logger.Default().Info("HttpServerStop stopped")
}
