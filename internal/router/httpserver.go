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
	HttpSrvHandler = &http.Server{
		Addr:           configs.GetServerConfig().Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(configs.GetServerConfig().ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(configs.GetServerConfig().WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << uint(configs.GetServerConfig().MaxHeaderBytes),
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
