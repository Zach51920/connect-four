package models

const (
	GameTypeBot     = "BOT"
	GameTypeLocal   = "LOCAL"
	GameTypeBotOnly = "BOT_ONLY"
)

type MakeMoveRequest struct {
	Column int `form:"column"`
}

type CreateGameRequest struct {
	Type string `form:"game_type"`
}

type BotConfigRequest struct {
	ID               string `form:"id"`
	Difficulty       int    `form:"difficulty"`
	MistakeFrequency int    `form:"mistake_frequency"`
	IsRandom         string `form:"is_random"`
}
