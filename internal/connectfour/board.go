package connectfour

import "errors"

const (
	DefaultBoardRows    = 6
	DefaultBoardColumns = 7

	winWeight    = 1000
	centerWeight = 5
)

var ErrInvalidColumn = errors.New("invalid column")

type Board struct {
	Cells        [][]rune
	winningCells [][2]int // Add this field to store winning cell coordinates
}

func NewBoard(rows, cols int) *Board {
	cells := make([][]rune, rows)
	for i := range cells {
		cells[i] = make([]rune, cols)
	}
	return &Board{Cells: cells}
}

func (b *Board) Copy() *Board {
	board := NewBoard(b.NumRows(), b.NumCols())
	for i, row := range b.Cells {
		copy(board.Cells[i], row)
	}
	return board
}

func (b *Board) Insert(token rune, column int) error {
	if b.NumRows() == 0 {
		return errors.New("board has no rows")
	}

	if column >= len(b.Cells[0]) {
		return ErrInvalidColumn
	}

	// insert the token in the first empty row in the column
	for i := b.NumRows() - 1; i >= 0; i-- {
		if b.Cells[i][column] == 0 {
			b.Cells[i][column] = token
			return nil
		}
	}
	return ErrInvalidColumn
}

func (b *Board) CheckWin(token rune) bool {
	b.winningCells = nil // Reset winning Cells before checking

	// check rows
	for i := 0; i < b.NumRows(); i++ {
		for j := 0; j <= b.NumCols()-4; j++ {
			if b.Cells[i][j] == token && b.Cells[i][j+1] == token && b.Cells[i][j+2] == token && b.Cells[i][j+3] == token {
				b.winningCells = [][2]int{{i, j}, {i, j + 1}, {i, j + 2}, {i, j + 3}}
				return true
			}
		}
	}
	// check columns
	for i := 0; i <= b.NumRows()-4; i++ {
		for j := 0; j < b.NumCols(); j++ {
			if b.Cells[i][j] == token && b.Cells[i+1][j] == token && b.Cells[i+2][j] == token && b.Cells[i+3][j] == token {
				b.winningCells = [][2]int{{i, j}, {i + 1, j}, {i + 2, j}, {i + 3, j}}
				return true
			}
		}
	}
	// check diagonals (top-left to bottom-right)
	for i := 0; i <= b.NumRows()-4; i++ {
		for j := 0; j <= b.NumCols()-4; j++ {
			if b.Cells[i][j] == token && b.Cells[i+1][j+1] == token && b.Cells[i+2][j+2] == token && b.Cells[i+3][j+3] == token {
				b.winningCells = [][2]int{{i, j}, {i + 1, j + 1}, {i + 2, j + 2}, {i + 3, j + 3}}
				return true
			}
		}
	}
	// check diagonals (bottom-left to top-right)
	for i := 3; i < b.NumRows(); i++ {
		for j := 0; j <= b.NumCols()-4; j++ {
			if b.Cells[i][j] == token && b.Cells[i-1][j+1] == token && b.Cells[i-2][j+2] == token && b.Cells[i-3][j+3] == token {
				b.winningCells = [][2]int{{i, j}, {i - 1, j + 1}, {i - 2, j + 2}, {i - 3, j + 3}}
				return true
			}
		}
	}
	//b.winningCells = nil
	return false
}

func (b *Board) NumRows() int {
	return len(b.Cells)
}

func (b *Board) NumCols() int {
	if len(b.Cells) == 0 {
		return 0
	}
	return len(b.Cells[0])
}

func (b *Board) GetCell(row, col int) rune {
	if row >= b.NumRows() || col >= b.NumCols() {
		return 0
	}
	return b.Cells[row][col]
}

func (b *Board) SetCell(row, col int, token rune) {
	b.Cells[row][col] = token
}

func (b *Board) IsColumnFull(col int) bool {
	if b.NumRows() == 0 || b.NumCols() == 0 {
		return true
	}

	if len(b.Cells[0]) <= col {
		return true
	}
	return b.Cells[0][col] != 0
}

func (b *Board) IsFull() bool {
	if b.NumRows() == 0 || b.NumCols() == 0 {
		return true
	}

	// check the first row for any open slots
	for _, col := range b.Cells[0] {
		if col == 0 {
			return false
		}
	}
	return true
}

func (b *Board) IsWinningCell(row, col int) bool {
	if b.winningCells == nil {
		return false
	}

	for _, cell := range b.winningCells {
		if cell[0] == row && cell[1] == col {
			return true
		}
	}
	return false
}

func (b *Board) Evaluate(token, opToken rune) float64 {
	var score float64

	// Check for win/loss
	if b.CheckWin(token) {
		return winWeight
	}
	if b.CheckWin(opToken) {
		return -winWeight
	}

	// Evaluate center control
	center := b.NumCols() / 2
	centerCount := 0
	for row := 0; row < b.NumRows(); row++ {
		if b.GetCell(row, center) == token {
			centerCount++
		}
	}
	score += float64(centerCount * centerWeight)

	// Evaluate potential threats
	score += b.evaluateThreats(token)
	score -= b.evaluateThreats(opToken)
	return score
}

func (b *Board) evaluateThreats(token rune) float64 {
	var score float64

	// Check horizontal threats
	for row := 0; row < b.NumRows(); row++ {
		for col := 0; col < b.NumCols()-3; col++ {
			window := []rune{
				b.GetCell(row, col),
				b.GetCell(row, col+1),
				b.GetCell(row, col+2),
				b.GetCell(row, col+3),
			}
			score += b.evaluateWindow(window, token)
		}
	}

	// Check vertical threats
	for row := 0; row < b.NumRows()-3; row++ {
		for col := 0; col < b.NumCols(); col++ {
			window := []rune{
				b.GetCell(row, col),
				b.GetCell(row+1, col),
				b.GetCell(row+2, col),
				b.GetCell(row+3, col),
			}
			score += b.evaluateWindow(window, token)
		}
	}

	// Check diagonal threats (positive slope)
	for row := 0; row < b.NumRows()-3; row++ {
		for col := 0; col < b.NumCols()-3; col++ {
			window := []rune{
				b.GetCell(row, col),
				b.GetCell(row+1, col+1),
				b.GetCell(row+2, col+2),
				b.GetCell(row+3, col+3),
			}
			score += b.evaluateWindow(window, token)
		}
	}

	// Check diagonal threats (negative slope)
	for row := 3; row < b.NumRows(); row++ {
		for col := 0; col < b.NumCols()-3; col++ {
			window := []rune{
				b.GetCell(row, col),
				b.GetCell(row-1, col+1),
				b.GetCell(row-2, col+2),
				b.GetCell(row-3, col+3),
			}
			score += b.evaluateWindow(window, token)
		}
	}

	return score
}

func (b *Board) evaluateWindow(window []rune, token rune) float64 {
	tokenCount := 0
	emptyCount := 0

	for _, piece := range window {
		if piece == token {
			tokenCount++
		} else if piece == 0 {
			emptyCount++
		}
	}

	if tokenCount == 3 && emptyCount == 1 {
		return 5
	} else if tokenCount == 2 && emptyCount == 2 {
		return 2
	}

	return 0
}

func (b *Board) validColumns() []int {
	board := *b
	validCols := make([]int, 0, board.NumCols())
	for col := 0; col < board.NumCols(); col++ {
		if !board.IsColumnFull(col) {
			validCols = append(validCols, col)
		}
	}
	return validCols
}
