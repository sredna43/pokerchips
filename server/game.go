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

func handleError(err *models.Error, res *models.Response) []byte {
	log.Println("handling error " + err.Cause)
	res.Error = err
	b, jsonError := json.Marshal(*res)
	if jsonError != nil {
		log.Println("error: ", err)
		return []byte("Unknown error occured... " + jsonError.Error())
	}
	return b
}

func updateGameState(e []byte) []byte {
	var req *models.Request
	var player *models.Player
	res := &models.Response{}
	err := json.Unmarshal(e, &req)
	if err != nil {
		return handleError(&models.Error{"json decoding error", err.Error()}, res)
	}
	if _, ok := gameState.Players[req.Player.Name]; !ok && req.Action != "new_player" {
		return handleError(&models.Error{"player not found", "Could not find player with name " + req.Player.Name + " in players list."}, res)
	} else if req.Action != "new_player" {
		player = gameState.Players[req.Player.Name]
	}
	switch req.Action {
	case "new_player":
		ok := true
		for _, val := range gameState.Players {
			if req.Player.Name == val.Name {
				ok = false
				return handleError(&models.Error{"duplicate player", "Player name already exists, please pick a new one"}, res)
			}
		}
		if ok {
			player = &models.Player{
				Spot:   len(gameState.Players),
				Name:   req.Player.Name,
				Chips:  lobby.Settings.InitialChips,
				Folded: false,
			}
			res.Message = "added player " + player.Name
		}
	case "remove_player":
		delete(gameState.Players, player.Name)
		for _, val := range gameState.Players {
			if val.Spot > req.Player.Spot {
				val.Spot -= 1
			}
		}
		res.Message = "removed player " + player.Name
		player = req.Player
	case "check":
		log.Println("check")
		res.Message = player.Name + " checks"
		gameState.Turn += 1
	case "bet":
		log.Println("bet")
		player.Chips -= req.Amount
		gameState.Pot += req.Amount
		res.Message = player.Name + " bets " + string(req.Amount)
		gameState.Turn += 1
	case "fold":
		log.Println("fold")
		player.Folded = true
		res.Message = player.Name + " folds"
		gameState.Turn += 1
	}
	if req.Action != "remove_player" {
		gameState.Players[player.Name] = player
		res.GameState = &gameState
	}
	res.Player = player
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
