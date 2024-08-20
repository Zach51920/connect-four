package main

import (
	"github.com/Zach51920/connect-four/internal/server"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"log/slog"
	"os"
)

func main() {
	// set slog to debug
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// create and run the server
	s := server.New(":8080")
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
