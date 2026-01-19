package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-api/auth"
	"github.com/example/go-api/controller"
	"github.com/example/go-api/internal/config"
	"github.com/example/go-api/middleware"
	"github.com/example/go-api/repository"
	"github.com/example/go-api/service"
)

// Server represents the HTTP server.
type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	http   *http.Server
	engine *gin.Engine
}

// New configures routing, handlers, and middleware.
func New(cfg *config.Config, logger *zap.Logger, db *sql.DB) (*Server, error) {
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(
		middleware.RequestID(),
		middleware.BodyLimit(cfg.Server.MaxBodyBytes),
		middleware.RateLimit(cfg.RateLimit.RPS),
		middleware.CORS(cfg.CORS),
		middleware.AccessLogger(logger),
		middleware.Recovery(logger),
	)

	userRepo := repository.NewUserRepository(db)
	jwtManager, err := auth.NewJWTManager(cfg.Auth)
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(db, userRepo)
	authService := service.NewAuthService(userRepo, jwtManager)

	userController := controller.NewUserController(userService, logger)
	authController := controller.NewAuthController(authService, logger)

	setupRoutes(engine, cfg, logger, userController, authController, jwtManager)

	return &Server{cfg: cfg, logger: logger, engine: engine}, nil
}

// Start boots the HTTP server and blocks until it exits.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Address, s.cfg.Server.Port)
	s.http = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}

	s.logger.Info("server starting", zap.String("addr", addr), zap.String("env", s.cfg.App.Env))
	err := s.http.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown attempts graceful HTTP server stop.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.http == nil {
		return nil
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	s.logger.Info("shutting down http server")
	return s.http.Shutdown(shutdownCtx)
}

// Engine exposes the gin engine for testing.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}
