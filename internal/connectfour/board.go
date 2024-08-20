package connectfour

import "errors"

const (
	DefaultBoardRows    = 6
	DefaultBoardColumns = 7

	winWeight    = 1000
	centerWeight = 5
)

var ErrInvalidColumn = errors.New("invalid column")

type Board [][]rune

func NewBoard(rows, cols int) *Board {
	cells := make(Board, rows)
	for i := range cells {
		cells[i] = make([]rune, cols)
	}
	return &cells
}

func (b *Board) Copy() *Board {
	board := NewBoard(b.NumRows(), b.NumCols())
	for i, row := range *b {
		for j, token := range row {
			board.SetCell(i, j, token)
		}
	}
	return board
}

func (b *Board) Insert(token rune, column int) error {
	if b.NumRows() == 0 {
		return errors.New("board has no rows")
	}

	board := *b
	if column >= len(board[0]) {
		return ErrInvalidColumn
	}

	// insert the token in the first empty row in the column
	for i := b.NumRows() - 1; i >= 0; i-- {
		row := board[i]
		if row[column] == 0 {
			row[column] = token
			return nil
		}
	}
	return ErrInvalidColumn
}

func (b *Board) CheckWin(token rune) bool {
	board := *b
	rows := b.NumRows()
	cols := b.NumCols()

	// check rows
	for i := 0; i < rows; i++ {
		for j := 0; j < cols-3; j++ {
			if board[i][j] == token && board[i][j+1] == token && board[i][j+2] == token && board[i][j+3] == token {
				return true
			}
		}
	}
	// check columns
	for i := 0; i < rows-3; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == token && board[i+1][j] == token && board[i+2][j] == token && board[i+3][j] == token {
				return true
			}
		}
	}
	// check diagonals (top-left to bottom-right)
	for i := 0; i < rows-3; i++ {
		for j := 0; j < cols-3; j++ {
			if board[i][j] == token && board[i+1][j+1] == token && board[i+2][j+2] == token && board[i+3][j+3] == token {
				return true
			}
		}
	}
	// check diagonals (bottom-left to top-right)
	for i := 3; i < rows; i++ {
		for j := 0; j < cols-3; j++ {
			if board[i][j] == token && board[i-1][j+1] == token && board[i-2][j+2] == token && board[i-3][j+3] == token {
				return true
			}
		}
	}
	return false
}

func (b *Board) NumRows() int {
	return len(*b)

}

func (b *Board) NumCols() int {
	if len(*b) == 0 {
		return 0
	}
	board := *b
	return len(board[0])
}

func (b *Board) GetCell(row, col int) rune {
	if row >= b.NumRows() || col >= b.NumCols() {
		return 0
	}
	board := *b
	return board[row][col]
}

func (b *Board) SetCell(row, col int, token rune) {
	board := *b
	board[row][col] = token
}

func (b *Board) IsColumnFull(col int) bool {
	if b.NumRows() == 0 || b.NumCols() == 0 {
		return true
	}

	board := *b
	if len(board[0]) <= col {
		return true
	}
	return board[0][col] != 0
}

func (b *Board) IsFull() bool {
	if b.NumRows() == 0 || b.NumCols() == 0 {
		return true
	}

	// check the first row for any open slots
	board := *b
	for _, col := range board[0] {
		if col == 0 {
			return false
		}
	}
	return true
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
