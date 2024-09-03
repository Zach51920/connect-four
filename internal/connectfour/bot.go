package connectfour

import (
	"fmt"
	"log/slog"
	"math/rand"
)

type Config struct {
	MistakeFrequency int
	Difficulty       int
	Randomize        bool
}

func DefaultConfig() *Config {
	return &Config{
		MistakeFrequency: 5,
		Difficulty:       6,
		Randomize:        true,
	}
}

func (c *Config) SetMistakeFrequency(freq int) *Config { c.MistakeFrequency = freq; return c }

func (c *Config) SetDifficulty(difficulty int) *Config { c.Difficulty = difficulty; return c }

func (c *Config) IncludeRandomization(randomize bool) *Config { c.Randomize = randomize; return c }

type Strategy interface {
	Name() string
	Suggest(board *Board, token rune) int
}

type BotPlayer struct {
	BasePlayer
	Config   *Config
	strategy Strategy
}

func (p *BotPlayer) Evaluate(board *Board) int {
	if col := p.initialEval(board); col != -1 {
		return col
	}
	return p.strategy.Suggest(board, p.token)
}

func (p *BotPlayer) initialEval(board *Board) int {
	if rand.Intn(100-p.Config.MistakeFrequency+1) == 0 {
		slog.Debug("bot is making an intentional mistake")
		// make a mistake, return random column
		validCols := board.validColumns()
		return validCols[rand.Intn(len(validCols))]
	}

	// check for immediate win or block
	if col, isWin := isWinningTurn(board, p.token); isWin {
		return col
	}
	opToken := tokenSwitch[p.token]
	if col, isWin := isWinningTurn(board, opToken); isWin {
		return col
	}
	return -1
}

func (p *BotPlayer) Strategy() string { return p.strategy.Name() }

func randomUsername() string {
	adjectives := []string{"Squeaky", "Fluffy", "Snazzy", "Clumsy", "Derpy", "Zesty", "Wacky"}
	nouns := []string{"Whale", "Pigeon", "Donut", "Panda", "Noodle", "Giraffe", "Raccoon"}

	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s %s", adj, noun)
}

func isWinningTurn(board *Board, token rune) (int, bool) {
	for _, col := range board.validColumns() {
		tmpBoard := board.Copy()
		_ = tmpBoard.Insert(token, col)
		if tmpBoard.CheckWin(token) {
			return col, true
		}
	}
	return -1, false
}
