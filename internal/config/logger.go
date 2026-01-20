package config

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger builds a zap logger configured for JSON output.
func NewLogger(level, filePath string) (*zap.Logger, error) {
	lvl := zapcore.InfoLevel
	if err := lvl.Set(strings.ToLower(level)); err != nil {
		lvl = zapcore.InfoLevel
	}

	if filePath == "" {
		filePath = "logs/app.log"
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return nil, err
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
		},
		OutputPaths:      []string{"stdout", filePath},
		ErrorOutputPaths: []string{"stderr", filePath},
	}

	return cfg.Build()
}
