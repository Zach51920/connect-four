package sessions

import (
	"fmt"
	"github.com/Zach51920/connect-four/internal/connectfour"
	views "github.com/Zach51920/connect-four/internal/views2"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strings"
	"time"
)

const StreamRefreshInterval = 10 * time.Second

type Session struct {
	ID       string
	Game     *connectfour.Game
	LastUsed time.Time

	refreshCh   chan bool
	shutdownCh  chan struct{}
	isStreaming bool
}

func (s *Session) SetGame(game *connectfour.Game) {
	s.Game = game
}

func (s *Session) Refresh() {
	s.refreshCh <- true
}

func (s *Session) CloseStream() {
	if s.isStreaming && s.shutdownCh != nil {
		close(s.shutdownCh)
	}
}

func (s *Session) Stream(c *gin.Context) {
	// check if we're already streaming
	if s.isStreaming {
		slog.Debug("Unable to start stream", "error", "client stream already exists")
		return
	}
	s.isStreaming = true
	s.refreshCh = make(chan bool)
	s.shutdownCh = make(chan struct{})
	defer func() { s.isStreaming = false }()

	// let the client know our intentions
	slog.Info("Starting SSE stream", "session_id", s.ID)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	closeCh := c.Writer.CloseNotify()

	// send an initial message to confirm we're connected
	if _, err := c.Writer.WriteString("event: connection\ndata: SSE connection established\n\n"); err != nil {
		slog.Error("Error establishing SSE connection", "error", err)
		return
	}
	slog.Debug("SSE connection established", "session_id", s.ID)

	ticker := time.NewTicker(StreamRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCh:
			slog.Debug("Stream shutdown triggered")
			return
		case <-closeCh:
			slog.Debug("Client closed connection", "session_id", s.ID)
			return
		case <-ticker.C:
			s.render(c)
		case <-s.refreshCh:
			s.render(c)
		}
	}
}

func (s *Session) render(c *gin.Context) {
	s.LastUsed = time.Now()
	slog.Debug("Refreshing game view", "session_id", s.ID)

	boardComponent := views.ConnectFourBoard(s.Game, *s.Game.Board)
	scoreComponent := views.ScoreCard(s.Game)

	boardHTML := new(strings.Builder)
	scoreHTML := new(strings.Builder)

	if err := boardComponent.Render(c.Request.Context(), boardHTML); err != nil {
		slog.Error("Failed to render board", "error", err)
		return
	}
	if err := scoreComponent.Render(c.Request.Context(), scoreHTML); err != nil {
		slog.Error("Failed to render score", "error", err)
		return
	}

	if _, err := fmt.Fprintf(c.Writer, "event: board-update\ndata: %s\n\n", boardHTML.String()); err != nil {
		slog.Error("Failed to write board-update", "error", err)
	}
	if _, err := fmt.Fprintf(c.Writer, "event: score-update\ndata: %s\n\n", scoreHTML.String()); err != nil {
		slog.Error("Failed to write score-update", "error", err)
	}
	c.Writer.Flush()
}
