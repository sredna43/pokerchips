package models

type Error struct {
	Player *Player `json:"player"`
	Cause   string `json:"cause"`
	Message string `json:"message"`
}
