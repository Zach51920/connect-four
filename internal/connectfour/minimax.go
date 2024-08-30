package connectfour

import (
	"math"
	"math/rand"
)

const (
	DefaultMinimaxDepth = 6
	RandomnessFactor    = 0.1
	MistakeFrequency    = 15 // every 15 * depth allow the bot to make a mistake
)

func NewMinimaxBot(token rune) *BotPlayer {
	return &BotPlayer{
		Strategy:   NewMinimaxStrat(DefaultMinimaxDepth),
		BasePlayer: NewBasePlayer(randomUsername(), token),
	}
}

type MinimaxStrat struct{ maxDepth int }

func NewMinimaxStrat(depth int) *MinimaxStrat {
	return &MinimaxStrat{maxDepth: depth}
}

func (p *BotPlayer) MakeBestMove(board *Board) {
	col := p.Strategy.Suggest(board, p.token)
	_ = p.MakeMove(board, col)
}

func (m *MinimaxStrat) Suggest(board *Board, token rune) int {
	bestCol := -1
	bestScore := math.Inf(-1)
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	opToken := tokenSwitch[token]

	// bots with a depth < 7 have a chance to make mistakes
	if m.maxDepth < 7 {
		freq := MistakeFrequency * m.maxDepth
		if rand.Intn(freq) == 0 {
			// make a mistake, return random column
			validCols := board.validColumns()
			col := validCols[rand.Intn(len(validCols))]
			return col
		}
	}

	// Check for immediate win or block
	if col, isWin := m.isWinningTurn(board, token); isWin {
		return col
	}
	if col, isWin := m.isWinningTurn(board, opToken); isWin {
		return col
	}

	for _, col := range board.validColumns() {
		tmpBoard := board.Copy().Insert(token, col)
		score := m.minimax(tmpBoard, token, opToken, m.maxDepth-1, false, alpha, beta)

		// Add randomness, smarter bots are less random
		randWeight := 1 - RandomnessFactor*float64(m.maxDepth)
		score += rand.Float64() * randWeight

		if score > bestScore {
			bestScore = score
			bestCol = col
		}

		alpha = math.Max(alpha, score)
		if beta <= alpha {
			break
		}
	}
	return bestCol
}

// SetSkill sets the maxDepth for the minimax algorithm
func (m *MinimaxStrat) SetSkill(depth int) { m.maxDepth = depth }

// Skill returns the max depth of the minimax algorithm
func (m *MinimaxStrat) Skill() int { return m.maxDepth }

func (m *MinimaxStrat) minimax(board *Board, token, opToken rune, depth int, isMaximizing bool, alpha, beta float64) float64 {
	if depth == 0 || board.IsFull() || board.CheckWin(token) || board.CheckWin(opToken) {
		return board.Evaluate(token, opToken)
	}
	validCols := board.validColumns()

	if isMaximizing {
		maxEval := math.Inf(-1)
		for _, col := range validCols {
			tmpBoard := board.Copy()
			_ = tmpBoard.Insert(token, col)
			eval := m.minimax(tmpBoard, token, opToken, depth-1, false, alpha, beta)
			maxEval = math.Max(maxEval, eval)
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := math.Inf(1)
		for _, col := range validCols {
			tmpBoard := board.Copy()
			_ = tmpBoard.Insert(opToken, col)
			eval := m.minimax(tmpBoard, token, opToken, depth-1, true, alpha, beta)
			minEval = math.Min(minEval, eval)
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

func (m *MinimaxStrat) isWinningTurn(board *Board, token rune) (int, bool) {
	for _, col := range board.validColumns() {
		tmpBoard := board.Copy()
		_ = tmpBoard.Insert(token, col)
		if tmpBoard.CheckWin(token) {
			return col, true
		}
	}
	return -1, false
}
