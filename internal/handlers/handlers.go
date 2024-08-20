package handlers

import (
	"fmt"
	"github.com/Zach51920/connect-four/internal/connectfour"
	"github.com/Zach51920/connect-four/internal/views"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	sessions *Sessions
}

func New() *Handler {
	return &Handler{sessions: NewSessions()}
}

func (h *Handler) Root(c *gin.Context) {
	render(c, views.Home())
}

func (h *Handler) GetGame(c *gin.Context) {
	game := h.getCurrentGame(c)
	if game == nil {
		renderHomeError(c, "You don't have any active games")
		return
	}
	render(c, views.Game(game))
}

func (h *Handler) CreateGame(c *gin.Context) {
	sessionID := c.GetString("session_id")
	if sessionID == "" {
		renderHomeError(c, "No active sessions found")
		return
	}

	var req CreateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to bind request", "error", err)
		renderHomeError(c, "An unexpected error occurred")
		return
	}

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
		renderHomeError(c, "Invalid game type")
		return
	}

	game := connectfour.NewGame(player1, player2)
	h.sessions.Set(sessionID, game)
	render(c, views.Game(game))
}

func (h *Handler) RestartGame(c *gin.Context) {
	game := h.getCurrentGame(c)
	if game == nil {
		renderHomeError(c, "You don't have any active games")
		return
	}
	game.Restart()

	// if the first move is made by a bot, make the move
	if bot, ok := game.Players[0].(*connectfour.BotPlayer); ok && game.HasHuman() {
		bot.MakeBestMove(game.Board)
	}
	render(c, views.Game(game))
}

func (h *Handler) SetDifficulty(c *gin.Context) {
	game := h.getCurrentGame(c)
	if game == nil {
		renderHomeError(c, "You don't have any active games")
		return
	}

	var req SetDifficultyRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to bind request", "error", err)
		renderHomeError(c, "An unexpected error occurred")
		return
	}
	slog.Debug("Setting bot difficulty", "difficulty", req.Difficulty, "id", req.ID)

	player, ok := game.GetPlayer(req.ID)
	if !ok {
		render(c, views.ErrorToast("Player not found"))
		return
	}
	bot, ok := player.(*connectfour.BotPlayer)
	if !ok {
		render(c, views.ErrorToast("Player is not a bot"))
		return
	}
	slog.Debug("Setting bot difficulty", "difficulty", req.Difficulty, "id", bot.ID())
	bot.Strategy.SetSkill(req.Difficulty)
	render(c, views.Game(game))
}

func (h *Handler) MakeMove(c *gin.Context) {
	var req MakeMoveRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to bind request", "error", err)
		renderHomeError(c, "An unexpected error occurred")
		return
	}
	game := h.getCurrentGame(c)
	if game == nil {
		renderHomeError(c, "You don't have any active games")
		return
	}
	defer render(c, views.Game(game)) // re-render the game

	// order players by turn
	// seems a bit jank but if the bot goes first their first move takes place in the Restart handler
	var players [2]connectfour.Player
	if game.IsTurn(game.Players[0]) {
		players = [2]connectfour.Player{game.Players[0], game.Players[1]}
	} else {
		players = [2]connectfour.Player{game.Players[1], game.Players[0]}
	}

	humanMoved := false
	for _, player := range players {
		human, ok := player.(*connectfour.HumanPlayer)
		if ok && !humanMoved {
			_ = human.MakeMove(game.Board, req.Column)
			state := game.State()
			if state == connectfour.GameStateWin {
				human.IncWins()
				render(c, views.Winner(human.Name()))
			} else if state == connectfour.GameStateDraw {
				render(c, views.Stalemate())
			} else {
				humanMoved = true
				continue
			}
			return
		}

		bot, ok := player.(*connectfour.BotPlayer)
		if ok {
			bot.MakeBestMove(game.Board)
			state := game.State()
			if state == connectfour.GameStateWin {
				render(c, views.Loser())
			} else if state == connectfour.GameStateDraw {
				render(c, views.Stalemate())
			} else {
				continue
			}
			return
		}
	}
}

func (h *Handler) getCurrentGame(c *gin.Context) *connectfour.Game {
	sessionID := c.GetString("session_id")
	if sessionID == "" {
		renderHomeError(c, "No active sessions found")
		return nil
	}
	return h.sessions.Get(sessionID)
}

func render(c *gin.Context, component templ.Component) {
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		_ = c.Error(fmt.Errorf("failed to render component: %w", err))
	}
}

func renderHomeError(c *gin.Context, message string) {
	render(c, views.ErrorToast(message))
	render(c, views.Home())
}
