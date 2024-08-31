package connectfour

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkMinimaxStrat_Suggest(b *testing.B) {
	// Create board states
	emptyBoard := NewBoard(DefaultBoardRows, DefaultBoardColumns)
	halfFullBoard := createHalfFullBoard()
	nearlyFullBoard := createNearlyFullBoard()

	// Test depths
	depths := []int{3, 5, 7}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("Depth_%d_EmptyBoard", depth), func(b *testing.B) {
			benchmarkSuggest(b, emptyBoard, depth)
		})
		b.Run(fmt.Sprintf("Depth_%d_HalfFullBoard", depth), func(b *testing.B) {
			benchmarkSuggest(b, halfFullBoard, depth)
		})
		b.Run(fmt.Sprintf("Depth_%d_NearlyFullBoard", depth), func(b *testing.B) {
			benchmarkSuggest(b, nearlyFullBoard, depth)
		})
	}
}

func benchmarkSuggest(b *testing.B, board *Board, depth int) {
	strat := NewMinimaxStrat(depth)
	token := 'X'

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strat.Suggest(board, token)
	}
}

func createHalfFullBoard() *Board {
	board := NewBoard(DefaultBoardRows, DefaultBoardColumns)
	token := 'X'

	for range DefaultBoardRows * DefaultBoardColumns / 2 {
		col := board.validColumns()[rand.Intn(len(board.validColumns()))]
		board.Insert(token, col)
		token = tokenSwitch[token]
	}
	return board
}

func createNearlyFullBoard() *Board {
	board := NewBoard(DefaultBoardRows, DefaultBoardColumns)
	token := 'X'

	for range DefaultBoardRows*DefaultBoardColumns - 5 {
		col := board.validColumns()[rand.Intn(len(board.validColumns()))]
		board.Insert(token, col)
		token = tokenSwitch[token]
	}
	return board
}
