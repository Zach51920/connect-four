package connectfour

import (
	"math"
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
	minBaseScore = 3.0
)

var tokenSwitch = map[rune]rune{'X': 'O', 'O': 'X'}

type Game struct {
	Players          [2]Player
	Board            *Board
	State            GameState
	Winner           Player
	MoveCount        int
	currentPlayerIdx int
}

func NewGame(player1, player2 Player) *Game {
	return &Game{
		Players: [2]Player{player1, player2},
		Board:   NewBoard(DefaultBoardRows, DefaultBoardColumns),
	}
}

func (g *Game) Restart() {
	g.State = GameStateNew
	g.Board = NewBoard(g.Board.NumRows(), g.Board.NumCols())
	g.currentPlayerIdx = 0
	g.MoveCount = 0
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
	_, isHuman := g.CurrentPlayer().(*HumanPlayer)
	return isHuman
}

func (g *Game) CurrentPlayer() Player {
	return g.Players[g.currentPlayerIdx]
}

func (g *Game) NextPlayer() Player {
	g.currentPlayerIdx = 1 - g.currentPlayerIdx // Toggle between 0 and 1
	return g.CurrentPlayer()
}

func (g *Game) IncMoveCount() { g.MoveCount++ }

func CalculateScore(player Player, board *Board) uint64 {
	// calculate the players score /100
	opToken := tokenSwitch[player.Token()]
	eval := board.Evaluate(player.Token(), opToken)
	clampedEval := math.Max(math.Min(eval, maxBaseScore), minBaseScore)

	// exponentially increase the score based on turn
	growthFactor := math.Pow(growthRate, float64(player.Turn()))
	score := clampedEval * growthFactor
	return uint64(score)
}
