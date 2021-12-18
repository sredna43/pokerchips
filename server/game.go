package main

import (
	"encoding/json"
	"log"

	"github.com/sredna43/pokerchips/models"
)

var gameState models.GameState = models.GameState{
	Playing: false,
	Players: make(map[string]*models.Player),
	Turn:    0,
	Dealer:  0,
}

var lobby models.Lobby = models.Lobby{
	Id:        "1a",
	GameState: &gameState,
	Settings: &models.Settings{
		InitialChips: 100,
	},
}

func updateGameState(e []byte) []byte {
	var req models.Request
	var errObject *models.Error
	var message string
	err := json.Unmarshal(e, &req)
	if err != nil {
		log.Println("error: ", err)
		errObject = &models.Error{"json decoding error", err.Error()}
	}
	if _, ok := gameState.Players[req.Player.Name]; !ok && req.Action != "new_player" {
		errObject = &models.Error{"player not found", "Could not find player with name " + req.Player.Name + " in players list."}
		req.Action = "error"
	}
	switch req.Action {
	case "check":
		message = req.Player.Name + " checks"
	case "bet":
		req.Player.Chips -= req.Amount
		gameState.Pot += req.Amount
		message = req.Player.Name + " bets " + string(req.Amount)
	case "fold":
		gameState.Players[req.Player.Name].Folded = true
		message = req.Player.Name + " folds"
	case "new_player":
		ok := true
		for _, val := range gameState.Players {
			if req.Player.Name == val.Name {
				ok = false
				errObject = &models.Error{"duplicate player", "Player name already exists, please pick a new one"}
			}
		}
		if ok {
			gameState.Players[req.Player.Name] = &models.Player{
				Spot:   len(gameState.Players),
				Name:   req.Player.Name,
				Chips:  lobby.Settings.InitialChips,
				Folded: false,
			}
			message = "added player " + req.Player.Name
		}
	case "remove_player":
		delete(gameState.Players, req.Player.Name)
		for _, val := range gameState.Players {
			if val.Spot > req.Player.Spot {
				val.Spot -= 1
			}
		}
		message = "removed player " + req.Player.Name
	}
	if gameState.Playing && gameState.Turn < len(gameState.Players) {
		gameState.Turn += 1
	}
	res := models.Response{
		Message:   message,
		Player:    gameState.Players[req.Player.Name],
		GameState: &gameState,
		Error:     errObject,
	}
	b, err := json.Marshal(res)
	if err != nil {
		log.Println("error: ", err)
	}
	return b
}

type Game struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newGame() *Game {
	return &Game{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (g *Game) run() {
	for {
		select {
		case client := <-g.register:
			g.clients[client] = true
		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.send)
			}
		case message := <-g.broadcast:
			s := updateGameState(message)
			for client := range g.clients {
				select {
				case client.send <- s:
				default:
					close(client.send)
					delete(g.clients, client)
				}
			}
		}
	}
}
