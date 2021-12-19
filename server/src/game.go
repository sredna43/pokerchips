package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/sredna43/pokerchips/models"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateLobby(n int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r1.Intn(len(letterBytes))]
	}
	return string(b)
}

func handleError(cause, message string, res *models.Response) []byte {
	log.Println("handling error " + cause)
	res.Error = &models.Error{
		Cause:   cause,
		Message: message,
	}
	res.Message = "error handled"
	b, jsonError := json.Marshal(*res)
	if jsonError != nil {
		log.Println("error: ", jsonError)
		return []byte("Unknown error occured... " + jsonError.Error())
	}
	return b
}

func updateGameState(req *models.Request, lobby *models.Lobby) []byte {
	gameState := lobby.GameState
	var player *models.Player
	res := &models.Response{}
	if _, ok := gameState.Players[req.Player.Name]; !ok && req.Action != "new_player" {
		return handleError("player not found", "Could not find player with name "+req.Player.Name+" in players list.", res)
	} else if req.Action != "new_player" {
		player = gameState.Players[req.Player.Name]
	} else if req.Player.Folded {
		return handleError("invalid action", "player has folded", res)
	}
	switch req.Action {
	case "new_player":
		log.Println("new_player")
		ok := true
		for _, val := range gameState.Players {
			if req.Player.Name == val.Name {
				ok = false
				return handleError("duplicate player", "player name already exists", res)
			}
		}
		if ok {
			player = &models.Player{
				Spot:   len(gameState.Players),
				Name:   req.Player.Name,
				Chips:  lobby.Settings.InitialChips,
				Folded: false,
			}
			if !lobby.Host {
				player.Host = true
			} else {
				player.Host = false
			}
			res.Message = "added player " + player.Name
		}
	case "remove_player":
		log.Println("remove_player")
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
		if req.Amount <= 0 {
			return handleError("invalid bet", "amount must be positive and non-zero", res)
		}
		if req.Amount > player.Chips {
			return handleError("invalid bet", "player does not have enough chips", res)
		}
		player.Chips -= req.Amount
		gameState.Pot += req.Amount
		res.Message = player.Name + " bets " + fmt.Sprint(req.Amount)
		gameState.Turn += 1
	case "fold":
		log.Println("fold")
		player.Folded = true
		res.Message = player.Name + " folds"
		gameState.Turn += 1
	case "new_hand":
		if !req.Player.Host {
			return handleError("action not allowed", "only the host can start a new hand", res)
		}
		log.Println("new_hand")
		for _, val := range gameState.Players {
			val.Folded = false
		}
		gameState.Pot = 0
		gameState.Dealer += 1
		if gameState.Dealer >= len(gameState.Players) {
			gameState.Dealer = 0
		}
		gameState.Turn = gameState.Dealer + 1
	case "start_game":
		if !req.Player.Host {
			return handleError("action not allowed", "only the host can start a new game", res)
		}
		gameState.Playing = true
	case "restart_game":
		if !req.Player.Host {
			return handleError("action not allowed", "only the host can start a new game", res)
		}
		gameState.Turn = 0
		gameState.Playing = true
		gameState.Dealer = 0
		gameState.Pot = 0
		for _, val := range gameState.Players {
			val.Chips = lobby.Settings.InitialChips
		}
	default:
		return handleError("invalid action", "unknown request action", res)
	}
	if gameState.Turn >= len(gameState.Players) {
		gameState.Turn = 0
	}
	if req.Action != "remove_player" {
		gameState.Players[player.Name] = player
		res.GameState = gameState
	}
	res.Player = player
	b, err := json.Marshal(res)
	if err != nil {
		log.Println("error: ", err)
	}
	return b
}

var lobbies = make(map[string]*models.Lobby)

func handleRequest(m []byte) []byte {
	var req *models.Request
	updateEligible := true
	err := json.Unmarshal(m, &req)
	res := &models.Response{}
	if err != nil {
		return handleError("json deconding error", err.Error(), res)
	}
	if _, ok := lobbies[req.Lobby]; !ok && req.Action != "new_game" {
		return handleError("invalid lobby id", "lobby id "+req.Lobby+" not found", res)
	}
	switch req.Action {
	case "new_game":
		lobbyId := generateLobby(3)
		lobbies[lobbyId] = &models.Lobby{
			GameState: &models.GameState{
				Players: make(map[string]*models.Player),
			},
			Settings: &models.Settings{},
		}
		res.Message = lobbyId
		updateEligible = false
	case "remove_game":
		delete(lobbies, req.Lobby)
		res.Message = "removed game " + req.Lobby
		updateEligible = false
	case "set_initial_chips":
		lobbies[req.Lobby].Settings.InitialChips = req.Amount
		res.Message = "set initial chips to " + fmt.Sprint(req.Amount)
		updateEligible = false
	case "set_max_players":
		lobbies[req.Lobby].Settings.MaxPlayers = req.Amount
		res.Message = "set max players to " + fmt.Sprint(req.Amount)
		updateEligible = false
	case "get_state":
		res.Message = "current state"
		res.GameState = lobbies[req.Lobby].GameState
		updateEligible = false
	}
	if updateEligible {
		return updateGameState(req, lobbies[req.Lobby])
	} else {
		b, err := json.Marshal(res)
		if err != nil {
			log.Println("error: ", err)
		}
		return b
	}
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
			s := handleRequest(message)
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