package repository

import (
	"context"
	"github.com/Zach51920/connect-four/internal/connectfour"
	"log/slog"
)

type Repository interface {
	SaveMove(ctx context.Context, game *connectfour.Game, player connectfour.Player, column int) error
}

type MockRepository struct{}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (r *MockRepository) SaveMove(ctx context.Context, game *connectfour.Game, player connectfour.Player, column int) error {
	slog.Debug("MOCK_REPO: save move")
	return nil
}
