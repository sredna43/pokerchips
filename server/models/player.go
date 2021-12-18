package models

type Player struct {
	Spot   int    `json:"spot"`
	Name   string `json:"name"`
	Chips  int    `json:"chips"`
	Folded bool   `json:"folded"`
}
