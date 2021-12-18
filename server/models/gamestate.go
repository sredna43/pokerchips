package models

type GameState struct {
	Playing bool               `json:"playing"`
	Players map[string]*Player `json:"players"`
	Turn    int                `json:"whose_turn"`
	Dealer  int                `json:"dealer"`
	Pot     int                `json:"pot"`
}
