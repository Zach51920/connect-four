package repository

import (
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
	ID        string    `bson:"_id,omitempty"`
	Player1   Player    `bson:"player1"`
	Player2   Player    `bson:"player2"`
	Moves     []Move    `bson:"moves"`
	Winner    *Player   `bson:"winner,omitempty"`
	MoveCount int       `bson:"move_count"`
	Timestamp time.Time `bson:"timestamp"`
}
