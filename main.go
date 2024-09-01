package main

import (
	"github.com/Zach51920/connect-four/internal/config"
	"github.com/Zach51920/connect-four/internal/server"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"log/slog"
	"os"
)

func main() {
	// load config
	cfg := config.Load(os.Getenv("CONFIG_PATH"))

	// initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.ParseLogLevel(),
	}))
	slog.SetDefault(logger)

	// create and run the server
	s := server.New(cfg.Server)
	defer s.Close()
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
