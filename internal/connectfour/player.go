package connectfour

import (
	"github.com/google/uuid"
	"log/slog"
	"math"
)

type Player interface {
	Name() string
	ID() string
	Token() rune
	Score() uint64
	AddScore(score uint64)
	MakeMove(board *Board, col int) error
	SetToken(token rune) *BasePlayer
	Reset() *BasePlayer
	Wins() int
	IncWins()
}

type BasePlayer struct {
	name  string
	token rune
	score uint64
	turn  int
	wins  int
	id    string
}

type HumanPlayer struct {
	BasePlayer
}

func NewBasePlayer(name string, token rune) BasePlayer {
	return BasePlayer{name: name, token: token, id: uuid.NewString()}
}

func NewHumanPlayer(name string, token rune) *HumanPlayer {
	return &HumanPlayer{BasePlayer: NewBasePlayer(name, token)}
}

func NewHumanPlayerPair() (*HumanPlayer, *HumanPlayer) {
	player1 := NewHumanPlayer("Player1", 'X')
	player2 := NewHumanPlayer("Player2", 'O')
	return player1, player2
}

func (p *BasePlayer) Name() string { return p.name }

func (p *BasePlayer) ID() string { return p.id }

func (p *BasePlayer) Score() uint64 { return p.score }

func (p *BasePlayer) AddScore(score uint64) { p.score += score }

func (p *BasePlayer) Wins() int { return p.wins }

func (p *BasePlayer) IncWins() { p.wins++ }

func (p *BasePlayer) Reset() *BasePlayer { p.score = 0; return p }

func (p *BasePlayer) MakeMove(board *Board, col int) error {
	if board.IsColumnFull(col) {
		return ErrInvalidMove
	}

	slog.Debug("making move", "col", col, "player", p.name)
	board.Insert(p.token, col)
	p.turn++

	// calculate the players score /100
	opToken := tokenSwitch[p.token]
	eval := board.Evaluate(p.token, opToken)
	clampedEval := math.Max(math.Min(eval, maxBaseScore), 0)

	// exponentially increase the score based on turn
	growthFactor := math.Pow(growthRate, float64(p.turn))
	score := clampedEval * growthFactor

	// update score
	slog.Debug("updating players score", "turn", p.turn, "player", p.name, "score", score, "growthFactor", growthFactor)
	p.AddScore(uint64(score))
	return nil
}

func (p *BasePlayer) SetToken(token rune) *BasePlayer { p.token = token; return p }

func (p *BasePlayer) Token() rune { return p.token }
