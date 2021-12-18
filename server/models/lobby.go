package models

type Lobby struct {
	Id        string     `json:"id"`
	GameState *GameState `json:"game_state"`
	Settings  *Settings  `json:"settings"`
}
