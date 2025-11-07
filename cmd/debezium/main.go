package main

import (
	"context"
	"debez/internal/config"
	v1 "debez/internal/transport/http/v1"
	"debez/pkg/logger"
	"debez/pkg/postgres"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./config/.env"
	}
	cfg, err := config.ParseConfig(envPath)
	if err != nil {
		log.Println("failed to parse config:", err)
		return
	}
	ctx, err = logger.New(ctx, cfg.Environment)
	if err != nil {
		log.Println("failed to initialize logger:", err)
		return
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Starting service...")
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Config:", zap.Any("config", cfg))

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to connect to db", zap.Error(err))
		return
	}

	server := v1.NewServer(cfg.Server.Port, db.Pool)
	err = server.RegisterHandler(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to register http handler", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.GetLoggerFromCtx(ctx).Info(ctx, "http server started", zap.Int("port", cfg.Server.Port))
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "http server error", zap.Error(err))
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	logger.GetLoggerFromCtx(ctx).Info(ctx, "shutting down service...")
	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.TimeOut)
	defer cancel()
	if err := server.Stop(shutdownCtx); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to stop http server", zap.Error(err))
	}
	db.Close()
	wg.Wait()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "service stopped")
}

/*
TODO
1) Дописать gRPC api
2) Подключить бд
3) Написать запросы
4) Кэш

*/
