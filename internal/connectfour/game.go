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
	GameStateStopped
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
	Winner  Player
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
	g.Winner = nil

	for _, player := range g.Players {
		player.Reset()
	}
}

func (g *Game) RefreshState() GameState {
	if g.State == GameStateCancelled || g.State == GameStateStopped {
		return g.State
	}

	g.State = GameStateOngoing
	if g.Board.IsFull() {
		g.State = GameStateDraw
	}
	for _, player := range g.Players {
		if g.Board.CheckWin(player.Token()) {
			g.State = GameStateWin
			g.Winner = player
			player.IncWins()
			break
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
	return g.State == GameStateNew || g.State == GameStateOngoing
}

func (g *Game) Stop() {
	g.State = GameStateStopped
}

func (g *Game) Resume() {
	g.State = GameStateNew
	g.RefreshState()
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
