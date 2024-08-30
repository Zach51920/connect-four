package connectfour

// TurnList represents the turns in the game
type TurnList struct {
	game               *Game
	currentPlayerIndex int
}

// NewTurnList creates a new TurnList for the given game
func NewTurnList(game *Game) *TurnList {
	return &TurnList{
		game:               game,
		currentPlayerIndex: 0,
	}
}

// Next moves to the next player and returns it
func (tl *TurnList) Next() Player {
	tl.currentPlayerIndex = 1 - tl.currentPlayerIndex // Toggle between 0 and 1
	return tl.Current()
}

// Current returns the current player
func (tl *TurnList) Current() Player {
	return tl.game.Players[tl.currentPlayerIndex]
}

// Reset sets the current player to the first player
func (tl *TurnList) Reset() {
	tl.currentPlayerIndex = 0
}
