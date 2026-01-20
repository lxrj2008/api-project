// Code generated API server entrypoint.
//
// @title Go API Boilerplate
// @version 1.0
// @description Simplified Go Web API boilerplate with JWT auth and SQL Server persistence.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"liangxiong/demo/internal/config"
	"liangxiong/demo/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(config.ConfigPath())
	if err != nil {
		panic(err)
	}

	logger, err := config.NewLogger(cfg.App.LogLevel, cfg.Logging.FilePath)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := config.NewDatabase(ctx, cfg.Database, logger)
	if err != nil {
		logger.Fatal("database init failed", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("db close failed", zap.Error(err))
		}
	}()

	srv, err := server.New(cfg, logger, db)
	if err != nil {
		logger.Fatal("server init failed", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("server stopped unexpectedly", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("signal received, shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}

	logger.Info("shutdown complete")
	os.Exit(0)
}
