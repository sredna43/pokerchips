package models

type Player struct {
	Host   bool   `json:"is_host"`
	Spot   int    `json:"spot"`
	Name   string `json:"name"`
	Chips  int    `json:"chips"`
	Folded bool   `json:"folded"`
}
