package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
)

func logMiddleware(c *gin.Context) {
	sessionID := c.GetString("session_id")
	slog.Info("received request", "path", c.Request.URL.Path, "session_id", sessionID)
	c.Next()
}

func sessionMiddleware(c *gin.Context) {
	s := sessions.Default(c)
	sessionID := s.Get("session_id")
	if sessionID == nil {
		sessionID = setSessionID(s)
	}
	c.Set("session_id", sessionID)
	c.Next()
}

func setSessionID(s sessions.Session) string {
	sessionID := uuid.New().String()
	slog.Debug("assigning session ID: not found in session", "session_id", sessionID)
	s.Set("session_id", sessionID)
	if err := s.Save(); err != nil {
		slog.Warn("failed to save session ID to cookie store", "error", err)
	}
	return sessionID
}
