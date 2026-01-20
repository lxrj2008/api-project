package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const defaultEnv = "dev"

// Config represents the full application configuration tree.
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Server    ServerConfig    `mapstructure:"server"`
	CORS      CORSConfig      `mapstructure:"cors"`
	Database  DatabaseConfig  `mapstructure:"db"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	RateLimit RateLimitConfig `mapstructure:"rateLimit"`
}

// AppConfig captures high-level application information.
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"logLevel"`
}

// ServerConfig controls HTTP server behaviour.
type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	MaxBodyBytes int64         `mapstructure:"maxBodyBytes"`
}

// CORSConfig enumerates allowed CORS options.
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
}

// DatabaseConfig contains SQL Server connection values.
type DatabaseConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

// AuthConfig glues JWT settings.
type AuthConfig struct {
	JWTSecret      string        `mapstructure:"jwtSecret"`
	Issuer         string        `mapstructure:"issuer"`
	Audience       string        `mapstructure:"audience"`
	AccessTokenTTL time.Duration `mapstructure:"accessTokenTTL"`
	Algorithm      string        `mapstructure:"algorithm"`
	PrivateKeyPath string        `mapstructure:"privateKeyPath"`
	PublicKeyPath  string        `mapstructure:"publicKeyPath"`
}

// LoggingConfig controls log output destinations.
type LoggingConfig struct {
	FilePath string `mapstructure:"filePath"`
}

// RateLimitConfig defines API throughput guardrails.
type RateLimitConfig struct {
	RPS int `mapstructure:"rps"`
}

// Load reads configuration for the current environment.
func Load(basePath string) (*Config, error) {
	cfg := &Config{}
	v := viper.New()

	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	if env == "" {
		env = defaultEnv
	}

	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath(basePath)
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate ensures the config contains all mandatory values.
func (c *Config) Validate() error {
	var missing []string

	if c.App.Name == "" {
		missing = append(missing, "app.name")
	}
	if c.App.Env == "" {
		missing = append(missing, "app.env")
	}
	if c.App.LogLevel == "" {
		missing = append(missing, "app.logLevel")
	}
	if c.Server.Address == "" || c.Server.Port == 0 {
		missing = append(missing, "server.address", "server.port")
	}
	if c.Server.MaxBodyBytes == 0 {
		missing = append(missing, "server.maxBodyBytes")
	}
	if c.Database.DSN == "" {
		missing = append(missing, "db.dsn")
	}
	if c.Database.ConnMaxLifetime == 0 {
		missing = append(missing, "db.connMaxLifetime")
	}
	if c.Auth.Issuer == "" {
		missing = append(missing, "auth.issuer")
	}
	if c.Auth.Audience == "" {
		missing = append(missing, "auth.audience")
	}
	if c.Auth.AccessTokenTTL == 0 {
		missing = append(missing, "auth.accessTokenTTL")
	}
	if c.Auth.Algorithm == "" {
		missing = append(missing, "auth.algorithm")
	}
	if c.Logging.FilePath == "" {
		missing = append(missing, "logging.filePath")
	}
	if c.RateLimit.RPS <= 0 {
		missing = append(missing, "rateLimit.rps")
	}

	switch strings.ToUpper(c.Auth.Algorithm) {
	case "HS256":
		if c.Auth.JWTSecret == "" {
			missing = append(missing, "auth.jwtSecret")
		}
	case "RS256":
		if c.Auth.PrivateKeyPath == "" || c.Auth.PublicKeyPath == "" {
			missing = append(missing, "auth.privateKeyPath", "auth.publicKeyPath")
		}
	default:
		missing = append(missing, "auth.algorithm")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}
	return nil
}

// ConfigPath resolves base path relative to executable.
func ConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return filepath.Clean(cwd)
}
