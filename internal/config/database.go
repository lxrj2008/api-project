package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

// NewDatabase creates a SQL Server connection pool.
func NewDatabase(ctx context.Context, cfg DatabaseConfig, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("database connection established")
	return db, nil
}
