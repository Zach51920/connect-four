package connectfour

import (
	"fmt"
	"log/slog"
	"math/rand"
)

const MistakeFrequency = 15 // every 15 * depth allow the bot to make a mistake

type Strategy interface {
	Suggest(board *Board, token rune) int
	SetDifficulty(skill int)
	Name() string
}

type BotPlayer struct {
	BasePlayer
	strategy   Strategy
	Difficulty int
}

func (p *BotPlayer) Evaluate(board *Board) int {
	if col := p.initialEval(board); col != -1 {
		return col
	}
	return p.strategy.Suggest(board, p.token)
}

func (p *BotPlayer) SetDifficulty(difficulty int) {
	p.Difficulty = difficulty
	p.strategy.SetDifficulty(difficulty)
}

func (p *BotPlayer) initialEval(board *Board) int {
	// bots with a difficulty < 7 have a chance to make mistakes
	if p.Difficulty < 7 {
		freq := MistakeFrequency * p.Difficulty
		if rand.Intn(freq) == 0 {
			slog.Debug("bot is making an intentional mistake")
			// make a mistake, return random column
			validCols := board.validColumns()
			return validCols[rand.Intn(len(validCols))]
		}
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
