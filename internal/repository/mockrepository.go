package repository

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"log/slog"
)

type MockRepository struct{}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (r *MockRepository) CreateGame(game *connectfour.Game) error {
	slog.Debug("MOCK_REPO: create game")
	return nil
}
