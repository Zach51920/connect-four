package connectfour

import (
	"log/slog"
	"math"
	"math/rand"
)

const (
	MinimaxRandomnessFactor = 0.1
	MinimaxDepthMultiplier  = 1
)

func NewMinimaxBot(token rune) *BotPlayer {
	config := DefaultConfig()
	return &BotPlayer{
		Config:     config,
		strategy:   NewMinimaxStrat(config),
		BasePlayer: NewBasePlayer(randomUsername(), token),
	}
}

type MinimaxStrat struct {
	Config *Config
}

func NewMinimaxStrat(config *Config) *MinimaxStrat {
	return &MinimaxStrat{Config: config}
}

func (m *MinimaxStrat) Suggest(board *Board, token rune) int {
	depth := m.Config.Difficulty * MinimaxDepthMultiplier
	slog.Debug("Suggesting move", "depth", depth, "randomize", m.Config.Randomize)

	bestCol := -1
	bestScore := math.Inf(-1)
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	opToken := tokenSwitch[token]

	for _, col := range board.validColumns() {
		tmpBoard := board.Copy().Insert(token, col)
		score := m.Minimax(tmpBoard, token, opToken, depth, false, alpha, beta)

		// Add randomness, smarter bots are less random
		if m.Config.Randomize {
			randWeight := 1 - MinimaxRandomnessFactor*float64(m.Config.Difficulty)
			score += rand.Float64() * randWeight
		}

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
