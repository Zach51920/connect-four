package repository

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{}
}

func (r *MongoRepository) CreateGame(game *connectfour.Game) error {
	return nil
}

func (r *MongoRepository) SaveMove(game *connectfour.Game, player connectfour.Player, column int) error {
	return nil
}
