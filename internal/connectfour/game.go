package connectfour

import (
	"time"
)

type GameState int

const (
	GameStateOngoing GameState = iota
	GameStateWin
	GameStateDraw
)

const (
	growthRate   = 1.02
	maxBaseScore = 100.0
)

var tokenSwitch = map[rune]rune{'X': 'O', 'O': 'X'}

type Meta struct {
	StartTime time.Time
	LastMove  time.Time
	NumMoves  int
}

type Game struct {
	Players    [2]Player
	Board      *Board
	Meta       *Meta
	currentIdx int
}

func NewGame(player1, player2 Player) *Game {
	return &Game{
		Players: [2]Player{player1, player2},
		Board:   NewBoard(DefaultBoardRows, DefaultBoardColumns),
		Meta:    &Meta{StartTime: time.Now(), LastMove: time.Now()},
	}
}

func (g *Game) Restart() {
	g.Board = NewBoard(g.Board.NumRows(), g.Board.NumCols())
	g.currentIdx = 0

	// switch the player order
	g.Players[0], g.Players[1] = g.Players[1], g.Players[0]

	// reset the players and set the first player as red
	g.Players[0].Reset().SetToken('X')
	g.Players[1].Reset().SetToken('O')
}

func (g *Game) State() GameState {
	if g.Board.IsFull() {
		return GameStateDraw
	}
	for _, player := range g.Players {
		if g.Board.CheckWin(player.Token()) {
			return GameStateWin
		}
	}
	return GameStateOngoing
}

func (g *Game) IsTurn(player Player) bool {
	return g.Players[g.currentIdx].Token() == player.Token()
}

func (g *Game) HasHuman() bool {
	return hasPlayerType[*HumanPlayer](g.Players)
}

func (g *Game) HasBot() bool {
	return hasPlayerType[*BotPlayer](g.Players)
}

func (g *Game) GetPlayer(playerID string) (Player, bool) {
	for _, p := range g.Players {
		if p.ID() == playerID {
			return p, true
		}
	}
	return nil, false
}

func hasPlayerType[T *BotPlayer | *HumanPlayer](players [2]Player) bool {
	for _, player := range players {
		if _, ok := player.(T); ok {
			return true
		}
	}
	return false
}
