package handlers

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"github.com/Zach51920/connect-four/internal/views2"
	"github.com/Zach51920/connect-four/sessions"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type Handlers struct {
	sessions sessions.Store
}

func New() *Handlers {
	return &Handlers{sessions: sessions.NewMemorySessionStore()}
}

func (h *Handlers) Home(c *gin.Context) {
	// if there's an active game, cancel it
	sessionID := c.GetString("session_id")
	sess, _ := h.sessions.Get(sessionID)
	if sess != nil && sess.Game != nil && sess.Game.InProgress() {
		render(c, views.WarningToast("The active game has been cancelled"))
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

	var req CreateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("Failed to bind CreateGameRequest", "error", err)
		return
	}

	// create the players according to the game type
	var player1, player2 connectfour.Player
	switch req.Type {
	case GameTypeBot:
		player1 = connectfour.NewHumanPlayer("Player 1", 'X')
		player2 = connectfour.NewMinimaxBot('O')
	case GameTypeLocal:
		player1, player2 = connectfour.NewHumanPlayerPair()
	case GameTypeBotOnly:
		player1 = connectfour.NewMinimaxBot('X')
		player2 = connectfour.NewMinimaxBot('O')
	default:
		h.handleCriticalErr(c, "Unknown game type")
		return
	}

	// create the game and assign it to the session
	game := connectfour.NewGame(player1, player2)
	sess.SetGame(game)

	// render the initial game board
	if err := views.Game(sess.Game).Render(c.Request.Context(), c.Writer); err != nil {
		h.handleCriticalErr(c, "Failed to render game board")
		return
	}
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

	// the main game loop
	game := sess.Game
	humanMoved := false
	for game.InProgress() {
		player := game.Turns.Current()

		if _, ok = player.(*connectfour.HumanPlayer); ok {
			// if human already moved just return because we're expecting more input
			if humanMoved {
				return
			}

			// make a move from the input
			var req MakeMoveRequest
			if err := c.ShouldBind(&req); err != nil {
				h.handleError(c, "An unexpected error has occurred")
				slog.Error("Failed to bind MakeMoveRequest", "error", err)
				return
			}
			if err := player.MakeMove(game.Board, req.Column); err != nil {
				h.handleError(c, "Invalid move selection")
				return
			}
			humanMoved = true
		} else if bot, ok := player.(*connectfour.BotPlayer); ok {
			// add some artificial delay
			timer := time.NewTimer(500 * time.Millisecond)
			bot.MakeBestMove(game.Board)
			<-timer.C
		}

		// refresh board and session
		game.RefreshState()
		sess.Refresh()
		if game.InProgress() {
			game.Turns.Next()
		}
	}
}

func (h *Handlers) SetDifficulty(c *gin.Context) {

}

func render(c *gin.Context, component templ.Component) {
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		slog.Error("Failed to render component", "error", err)
	}
}

func (h *Handlers) handleCriticalErr(c *gin.Context, message string) {
	h.Home(c) // send em home
	render(c, views.ErrorToast(message))
}

func (h *Handlers) handleError(c *gin.Context, message string) {
	render(c, views.ErrorToast(message))
}
