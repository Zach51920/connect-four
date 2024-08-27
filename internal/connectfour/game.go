package connectfour

import (
	"time"
)

type GameState int

const (
	GameStateNew GameState = iota
	GameStateOngoing
	GameStateWin
	GameStateDraw
	GameStateCancelled
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
	Players [2]Player
	Turns   *TurnList
	Board   *Board
	Meta    *Meta
	State   GameState
}

func NewGame(player1, player2 Player) *Game {
	game := &Game{
		Players: [2]Player{player1, player2},
		Board:   NewBoard(DefaultBoardRows, DefaultBoardColumns),
		Meta:    &Meta{StartTime: time.Now(), LastMove: time.Now()},
	}
	game.Turns = NewTurnList(game)
	return game
}

func (g *Game) Restart() {
	g.State = GameStateNew
	g.Board = NewBoard(g.Board.NumRows(), g.Board.NumCols())
	g.Turns.Reset()

	for _, player := range g.Players {
		player.Reset()
	}
}

func (g *Game) RefreshState() GameState {
	if g.State == GameStateCancelled {
		return GameStateCancelled
	}

	g.State = GameStateOngoing
	if g.Board.IsFull() {
		g.State = GameStateDraw
	}
	for _, player := range g.Players {
		if g.Board.CheckWin(player.Token()) {
			g.State = GameStateWin
		}
	}
	return g.State
}

func (g *Game) HasHuman() bool {
	for _, player := range g.Players {
		if _, isHuman := player.(*HumanPlayer); isHuman {
			return true
		}
	}
	return false
}

func (g *Game) InProgress() bool {
	return g.State != GameStateDraw && g.State != GameStateWin && g.State != GameStateCancelled
}

func (g *Game) Cancel() {
	g.State = GameStateCancelled
}

func (g *Game) ExpectHumanInput() bool {
	if g.State == GameStateDraw || g.State == GameStateWin {
		return false
	}
	_, isHuman := g.Turns.Current().(*HumanPlayer)
	return isHuman
}
