package models

type Request struct {
	Player *Player `json:"player"`
	Action string  `json:"action"`
	Amount int     `json:"amount"`
}
