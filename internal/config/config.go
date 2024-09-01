package config

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	LogLevel string        `yaml:"log_level"`
	Server   *ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Address     string `yaml:"address"`
	GinMode     string `yaml:"gin_mode"`
	WithMongoDB bool   `yaml:"with_mongodb"`
}

func Load(path string) *Config {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	config := new(Config)
	if err = yaml.Unmarshal(b, config); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
	return config
}

func (c *Config) ParseLogLevel() slog.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (c *ServerConfig) ParseGinMode() string {
	switch strings.ToLower(c.GinMode) {
	case "debug":
		return gin.DebugMode
	case "release":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.ReleaseMode
	}
}
