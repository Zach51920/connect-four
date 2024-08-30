package connectfour

import "errors"

const (
	DefaultBoardRows    = 6
	DefaultBoardColumns = 7
	WinLength           = 4

	winWeight    = 1000
	centerWeight = 5
)

var ErrInvalidMove = errors.New("invalid move")

type Board struct {
	Cells        [][]rune
	heights      []int
	lastMove     [2]int
	winningCells [][2]int
}

func NewBoard(rows, cols int) *Board {
	cells := make([][]rune, rows)
	for i := range cells {
		cells[i] = make([]rune, cols)
	}
	return &Board{
		Cells:   cells,
		heights: make([]int, cols),
	}
}

func (b *Board) Copy() *Board {
	newBoard := &Board{
		Cells:        make([][]rune, len(b.Cells)),
		heights:      make([]int, len(b.heights)),
		lastMove:     b.lastMove,
		winningCells: make([][2]int, len(b.winningCells)),
	}

	for i, row := range b.Cells {
		newBoard.Cells[i] = make([]rune, len(row))
		copy(newBoard.Cells[i], row)
	}
	copy(newBoard.heights, b.heights)
	copy(newBoard.winningCells, b.winningCells)

	return newBoard
}

func (b *Board) Insert(token rune, col int) error {
	if col < 0 || col >= len(b.Cells[0]) {
		return ErrInvalidMove
	}

	if b.heights[col] >= len(b.Cells) {
		return ErrInvalidMove
	}

	row := len(b.Cells) - 1 - b.heights[col]
	b.Cells[row][col] = token
	b.heights[col]++
	b.lastMove = [2]int{row, col}
	b.winningCells = nil // reset winning cells as the board state has changed

	return nil
}

func (b *Board) CheckWin(token rune) bool {
	row, col := b.lastMove[0], b.lastMove[1]

	// If the last move wasn't made by the token we're checking, return false immediately
	if b.Cells[row][col] != token {
		return false
	}

	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, dir := range directions {
		count := 1
		winningCells := [][2]int{{row, col}}

		// check in positive direction
		for i := 1; i < WinLength; i++ {
			r, c := row+i*dir[0], col+i*dir[1]
			if !b.isValidCell(r, c) || b.Cells[r][c] != token {
				break
			}
			count++
			winningCells = append(winningCells, [2]int{r, c})
		}

		// check in negative direction
		for i := 1; i < WinLength; i++ {
			r, c := row-i*dir[0], col-i*dir[1]
			if !b.isValidCell(r, c) || b.Cells[r][c] != token {
				break
			}
			count++
			winningCells = append(winningCells, [2]int{r, c})
		}

		if count >= WinLength {
			b.winningCells = winningCells
			return true
		}
	}

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
	for _, height := range b.heights {
		if height < len(b.Cells) {
			return false
		}
	}
	return true
}

func (b *Board) IsWinningCell(row, col int) bool {
	// if there are no winning cells, return false immediately
	if len(b.winningCells) == 0 {
		return false
	}

	// check if the given cell is in the winning cells slice
	for _, cell := range b.winningCells {
		if cell[0] == row && cell[1] == col {
			return true
		}
	}

	return false
}

func (b *Board) Evaluate(token, opToken rune) float64 {
	var score float64

	// check for win/loss
	if b.CheckWin(token) {
		return winWeight
	}
	if b.CheckWin(opToken) {
		return -winWeight
	}

	// evaluate center control
	center := b.NumCols() / 2
	centerCount := 0
	for row := 0; row < b.NumRows(); row++ {
		if b.GetCell(row, center) == token {
			centerCount++
		}
	}
	score += float64(centerCount * centerWeight)

	// evaluate potential threats
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

func (b *Board) isValidCell(row, col int) bool {
	return row >= 0 && row < len(b.Cells) && col >= 0 && col < len(b.Cells[0])
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
