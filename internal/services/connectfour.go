package services

import (
	"context"
	"errors"
	"github.com/Zach51920/connect-four/internal/connectfour"
	"github.com/Zach51920/connect-four/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type GameService struct {
	db *mongo.Database
}

func NewGameService(db *mongo.Database) *GameService {
	return &GameService{db: db}
}

func (s *GameService) CreateGame(ctx context.Context, sessionID string, req models.CreateGameRequest) (*connectfour.Game, error) {
	// create the players according to the game type
	var player1, player2 connectfour.Player
	switch req.Type {
	case models.GameTypeBot:
		player1 = connectfour.NewHumanPlayer("Player 1", 'X')
		player2 = connectfour.NewMinimaxBot('O')
	case models.GameTypeLocal:
		player1, player2 = connectfour.NewHumanPlayerPair()
	case models.GameTypeBotOnly:
		player1 = connectfour.NewMinimaxBot('X')
		player2 = connectfour.NewMinimaxBot('O')
	default:
		return nil, errors.New("unknown game type")
	}

	// create the game and assign it to the session
	return connectfour.NewGame(player1, player2), nil
}

func (s *GameService) SetDifficulty(ctx context.Context, players [2]connectfour.Player, req models.SetDifficultyRequest) error {
	for _, player := range players {
		if player.ID() != req.ID {
			continue
		}
		bot, ok := player.(*connectfour.BotPlayer)
		if !ok {
			return errors.New("invalid player type")
		}
		slog.Debug("Setting bot difficulty", "difficulty", req.Difficulty, "bot", bot.ID())
		bot.SetDifficulty(req.Difficulty)
		break
	}
	return nil
}

func (s *GameService) MakeMove(ctx context.Context, player connectfour.Player, game *connectfour.Game, col int) error {
	if player != game.Turns.Current() {
		return errors.New("not players turn")
	}

	// make the players move
	if err := player.MakeMove(game.Board, col); err != nil {
		return err
	}
	game.RefreshState()

	// todo: add move to db
	return nil
}
