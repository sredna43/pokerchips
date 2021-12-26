package models

type Response struct {
	Lobby	  string 	 `json:"lobby"`
	Message   string     `json:"message"`
	Player    *Player    `json:"player"`
	GameState *GameState `json:"game_state"`
	Error     *Error     `json:"error"`
}
