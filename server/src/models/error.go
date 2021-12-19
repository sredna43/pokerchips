package models

type Error struct {
	Cause   string `json:"cause"`
	Message string `json:"message"`
}
