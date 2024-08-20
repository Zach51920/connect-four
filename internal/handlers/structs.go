package handlers

const (
	GameTypeBot     = "BOT"
	GameTypeLocal   = "LOCAL"
	GameTypeBotOnly = "BOT_ONLY"
)

type CreateGameRequest struct {
	Type string `form:"game_type"`
}

type MakeMoveRequest struct {
	Column int `form:"column"`
}

type SetDifficultyRequest struct {
	Difficulty int    `form:"difficulty"`
	ID         string `form:"id"`
}
