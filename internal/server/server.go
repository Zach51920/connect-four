package server

import (
	"github.com/Zach51920/connect-four/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"os"
)

type Server struct {
	addr   string
	router *gin.Engine
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) init() error {
	r := gin.Default()

	// init cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "Content-Type")
	r.Use(cors.New(corsConfig))

	// init cookie store
	secret := os.Getenv("COOKIE_SECRET")
	store := cookie.NewStore([]byte(secret))
	session := sessions.Sessions("connect_four", store)
	r.Use(session)

	// register middlewares
	r.Use(sessionMiddleware)
	r.Use(logMiddleware)

	// register handlers
	handle := handlers.New()
	r.GET("/", handle.Home)
	r.POST("/game", handle.CreateGame)
	r.GET("/game/stream", handle.StreamGame)
	r.POST("/game/move", handle.MakeMove)
	r.POST("/game/difficulty", handle.SetDifficulty)
	r.POST("/game/restart", handle.RestartGame)

	s.router = r
	return nil
}

func (s *Server) Run() error {
	if err := s.init(); err != nil {
		return err
	}
	return s.router.Run(s.addr)
}
