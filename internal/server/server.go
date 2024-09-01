package server

import (
	"fmt"
	"github.com/Zach51920/connect-four/internal/config"
	"github.com/Zach51920/connect-four/internal/handlers"
	"github.com/Zach51920/connect-four/internal/mongo"
	"github.com/Zach51920/connect-four/internal/repository"
	"github.com/Zach51920/connect-four/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

type Server struct {
	router *gin.Engine
	config *config.ServerConfig
}

func New(cfg *config.ServerConfig) *Server {
	return &Server{config: cfg}
}

func (s *Server) init() error {
	// initialize gin router
	gin.SetMode(s.config.ParseGinMode())
	r := gin.New()
	r.Use(gin.Recovery())

	// init cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "Content-Type")
	r.Use(cors.New(corsConfig))

	// serve static files
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	publicDir := filepath.Join(wd, "public")
	r.Static("/public", publicDir)

	// init cookie store
	secret := os.Getenv("COOKIE_SECRET")
	store := cookie.NewStore([]byte(secret))
	session := sessions.Sessions("connect_four", store)
	r.Use(session)

	// register middleware
	r.Use(sessionMiddleware)
	r.Use(logMiddleware)

	// create dependencies
	provider, err := mongo.NewProvider(mongo.FromEnv())
	if err != nil {
		return fmt.Errorf("failed to initialize mongodb provider: %w", err)
	}
	if err = provider.Ping(); err != nil {
		return fmt.Errorf("failed to ping mongodb provider: %w", err)
	}
	repo := repository.NewMongoRepository(provider.DB())
	service := services.NewGameService(repo)
	handle := handlers.New(service)

	// register handlers
	r.GET("/", handle.Home)
	r.GET("/game", handle.GetGame)
	r.POST("/game", handle.CreateGame)
	r.GET("/game/stream", handle.StreamGame)
	r.POST("/game/move", handle.MakeMove)
	r.POST("/game/difficulty", handle.SetDifficulty)
	r.POST("/game/restart", handle.RestartGame)
	r.POST("/game/stop", handle.StopGame)
	r.GET("/settings", handle.Settings)

	s.router = r
	return nil
}

func (s *Server) Run() error {
	if err := s.init(); err != nil {
		return err
	}
	slog.Info("Starting server...")
	return s.router.Run(s.config.Address)
}
