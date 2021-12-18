package models

type Error struct {
	Cause   string `json:"error"`
	Message string `json:"message"`
}
