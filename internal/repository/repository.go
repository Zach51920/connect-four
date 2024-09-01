package repository

import "github.com/Zach51920/connect-four/internal/connectfour"

type Repository interface {
	CreateGame(game *connectfour.Game) error
	SaveMove(game *connectfour.Game, player connectfour.Player, column int) error
}
