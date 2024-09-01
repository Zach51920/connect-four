package connectfour

import (
	"github.com/google/uuid"
)

type Player interface {
	Name() string
	ID() string
	Token() rune
	SetToken(token rune) *BasePlayer
	Score() uint64
	AddScore(score uint64)
	Wins() int
	IncWins()
	Turn() int
	IncTurn()
	Reset() *BasePlayer
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

func (p *BasePlayer) AddScore(score uint64) { p.score += uint64(score) }

func (p *BasePlayer) Wins() int { return p.wins }

func (p *BasePlayer) IncWins() { p.wins++ }

func (p *BasePlayer) Reset() *BasePlayer { p.score = 0; return p }

func (p *BasePlayer) IncTurn() { p.turn++ }

func (p *BasePlayer) Turn() int { return p.turn }

func (p *BasePlayer) SetToken(token rune) *BasePlayer { p.token = token; return p }

func (p *BasePlayer) Token() rune { return p.token }
