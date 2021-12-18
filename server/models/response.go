package models

type Response struct {
	Message   string     `json:"message"`
	Player    *Player    `json:"player"`
	GameState *GameState `json:"game_state"`
	Error     *Error     `json:"error"`
}
