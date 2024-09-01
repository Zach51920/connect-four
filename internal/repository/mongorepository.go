package repository

import (
	"context"
	"time"

	"github.com/Zach51920/connect-four/internal/connectfour"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{collection: db.Collection("games")}
}

func (r *MongoRepository) SaveMove(ctx context.Context, game *connectfour.Game, player connectfour.Player, column int) error {
	mongoCtx, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	move := Move{
		ID:       game.MoveCount,
		Column:   column,
		PlayerID: player.ID(),
	}

	update := bson.M{
		"$setOnInsert": bson.M{
			"_id":       game.ID,
			"timestamp": time.Now(),
		},
		"$push": bson.M{"moves": move},
		"$inc":  bson.M{"move_count": 1},
		"$set": bson.M{
			"player1": mapPlayer(game.Players[0]),
			"player2": mapPlayer(game.Players[1]),
		},
	}

	if game.Winner != nil {
		update["$set"].(bson.M)["winner"] = game.Winner.ID()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(mongoCtx, bson.M{"_id": game.ID}, update, opts)
	return err
}

func mapPlayer(player connectfour.Player) Player {
	return Player{
		ID:       player.ID(),
		Strategy: player.Strategy(),
		Token:    player.Token(),
		Score:    player.Score(),
	}
}
