package handlers

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"github.com/Zach51920/connect-four/internal/models"
	"github.com/Zach51920/connect-four/internal/services"
	"github.com/Zach51920/connect-four/internal/sessions"
	"github.com/Zach51920/connect-four/internal/views"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type Handlers struct {
	sessions sessions.Store
	service  *services.GameService
}

func New(service *services.GameService) *Handlers {
	return &Handlers{
		sessions: sessions.NewMemorySessionStore(),
		service:  service,
	}
}

func (h *Handlers) Home(c *gin.Context) {
	// if there's an active game, cancel it
	sessionID := c.GetString("session_id")
	sess, _ := h.sessions.Get(sessionID)
	if sess != nil && sess.Game != nil && sess.Game.InProgress() {
		render(c, views.WarningToast("The active game has been aborted"))
		sess.Game.Cancel()
	}
	// render the home page
	render(c, views.Home())
}

func (h *Handlers) CreateGame(c *gin.Context) {
	// check if we have an active session
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil {
		sess = h.sessions.New(sessionID, nil)
	}

	// parse the request
	var req models.CreateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("Failed to bind CreateGameRequest", "error", err)
		return
	}

	// create the game and add assign it to our session
	game, err := h.service.CreateGame(req)
	if err != nil {
		h.handleCriticalErr(c, "Failed to create game")
		return
	}
	sess.SetGame(game)

	// render the initial game board
	if err = views.Game(sess.Game).Render(c.Request.Context(), c.Writer); err != nil {
		h.handleCriticalErr(c, "Failed to render game board")
		return
	}
}

func (h *Handlers) GetGame(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}
	sess.CloseStream()
	render(c, views.Game(sess.Game))
}

func (h *Handlers) StreamGame(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}
	sess.Stream(c)
}

func (h *Handlers) RestartGame(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}
	sess.Game.Restart()
	sess.Refresh()
}

func (h *Handlers) MakeMove(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}

	game := sess.Game
	game.Resume() // if the game was paused, playing a move should automatically resume game

	// the main game loop
	humanMoved := false
	for game.InProgress() {
		player := game.CurrentPlayer()

		if _, ok = player.(*connectfour.HumanPlayer); ok {
			// if human already moved just return because we're expecting more input
			if humanMoved {
				return
			}

			// make a move from the input
			var req models.MakeMoveRequest
			if err := c.ShouldBind(&req); err != nil {
				h.handleError(c, "An unexpected error has occurred")
				slog.Error("Failed to bind MakeMoveRequest", "error", err)
				return
			}
			if err := h.service.MakeMove(c, player, game, req.Column); err != nil {
				h.handleError(c, "Invalid move selection")
				return
			}
			humanMoved = true
		} else if bot, ok := player.(*connectfour.BotPlayer); ok {
			// add some artificial delay
			timer := time.NewTimer(300 * time.Millisecond)
			col := bot.Evaluate(game.Board)
			if err := h.service.MakeMove(c, player, game, col); err != nil {
				h.handleError(c, "Invalid move selection")
				return
			}
			<-timer.C
		}

		game.NextPlayer()
		sess.Refresh()
	}
}

func (h *Handlers) ConfigureBot(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}

	var req models.BotConfigRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to bind request", "error", err)
		h.handleError(c, "An unexpected error occurred")
		return
	}

	if err := h.service.UpdateBotConfig(sess.Game.Players, req); err != nil {
		h.handleError(c, "Failed to set Difficulty")
		return
	}
	render(c, views.SettingsModal(sess.Game))
}

func (h *Handlers) StopGame(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}
	sess.Game.Stop()
}

func (h *Handlers) Settings(c *gin.Context) {
	sessionID := c.GetString("session_id")
	sess, ok := h.sessions.Get(sessionID)
	if !ok || sess == nil || sess.Game == nil {
		h.handleCriticalErr(c, "Failed to get active game")
		return
	}
	render(c, views.SettingsModal(sess.Game))
}

func render(c *gin.Context, component templ.Component) {
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		slog.Error("Failed to render component", "error", err)
		_ = views.ErrorToast("An unexpected error has occurred").Render(c.Request.Context(), c.Writer)
	}
}

func (h *Handlers) handleCriticalErr(c *gin.Context, message string) {
	h.Home(c) // send em home
	render(c, views.ErrorToast(message))
}

func (h *Handlers) handleError(c *gin.Context, message string) {
	render(c, views.ErrorToast(message))
}
