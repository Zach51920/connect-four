package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Player struct {
	ID       string `bson:"id"`
	Strategy string `bson:"strategy"`
	Score    uint64 `bson:"score"`
	Token    rune   `bson:"token"`
}

type Move struct {
	ID       int    `bson:"id"`
	Column   int    `bson:"column"`
	PlayerID string `bson:"player_id"`
}

type Game struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Players   [2]Player          `bson:"players"`
	Moves     []Move             `bson:"moves"`
	Winner    *Player            `bson:"winner,omitempty"`
	MoveCount int                `bson:"move_count"`
	Timestamp time.Time          `bson:"timestamp"`
}
