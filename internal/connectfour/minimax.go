package connectfour

import (
	"math"
	"math/rand"
)

const (
	MinimaxDefaultDepth     = 6
	MinimaxRandomnessFactor = 0.1
)

func NewMinimaxBot(token rune) *BotPlayer {
	return &BotPlayer{
		Difficulty: MinimaxDefaultDepth,
		strategy:   NewMinimaxStrat(MinimaxDefaultDepth),
		BasePlayer: NewBasePlayer(randomUsername(), token),
	}
}

type MinimaxStrat struct{ maxDepth int }

func NewMinimaxStrat(depth int) *MinimaxStrat {
	return &MinimaxStrat{maxDepth: depth}
}

func (m *MinimaxStrat) Suggest(board *Board, token rune) int {
	bestCol := -1
	bestScore := math.Inf(-1)
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	opToken := tokenSwitch[token]

	for _, col := range board.validColumns() {
		tmpBoard := board.Copy().Insert(token, col)
		score := m.Minimax(tmpBoard, token, opToken, m.maxDepth-1, false, alpha, beta)

		// Add randomness, smarter bots are less random
		randWeight := 1 - MinimaxRandomnessFactor*float64(m.maxDepth)
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

func (m *MinimaxStrat) Minimax(board *Board, token, opToken rune, depth int, isMaximizing bool, alpha, beta float64) float64 {
	if depth == 0 || board.IsFull() || board.CheckWin(token) || board.CheckWin(opToken) {
		return board.Evaluate(token, opToken)
	}
	validCols := board.validColumns()

	if isMaximizing {
		maxEval := math.Inf(-1)
		for _, col := range validCols {
			tmpBoard := board.Copy()
			_ = tmpBoard.Insert(token, col)
			eval := m.Minimax(tmpBoard, token, opToken, depth-1, false, alpha, beta)
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
			eval := m.Minimax(tmpBoard, token, opToken, depth-1, true, alpha, beta)
			minEval = math.Min(minEval, eval)
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

func (m *MinimaxStrat) Name() string {
	return "MINMAX"
}

func (m *MinimaxStrat) SetDifficulty(depth int) { m.maxDepth = depth }
