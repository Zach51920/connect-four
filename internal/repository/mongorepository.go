package repository

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{}
}

func (r *MongoRepository) CreateGame(game *connectfour.Game) error {
	slog.Debug("MOCK_REPO: create game")
	return nil
}
