package models

type Lobby struct {
	Host      bool       `json:"has_host"`
	GameState *GameState `json:"game_state"`
	Settings  *Settings  `json:"settings"`
	CreatedOn  int64     `json:"created_on"`
}
