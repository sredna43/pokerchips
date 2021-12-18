package main

import "log"

const (
	bet = "0"
	fold = "1"
)

type Player struct {
	name string `json:"name"`
	chips int `json:"chips"`
}

type GameState struct {
	playing bool `json:"playing"`
	players map[int]Player `json:"players"`
	turn int `json:"whose_turn"`
	dealer int `json:"dealer"`
}

type Event struct {
	id int `json:"id"`
	action int `json:"action"`
	amount int `json:"amount"`
}

type Game struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

var gameState GameState = GameState{playing: false, players: make(map[int]Player), turn: 0, dealer: 0}

func newGame() *Game {
	return &Game{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func updateGameState(event string) []byte {
	if event == bet {
		return []byte("bet")
	}
	if event == fold {
		return []byte("fold")
	}
	if event == "np" {
		gameState.players[0] = Player{ name: "Player 1", chips: 100}
		log.Print(gameState)
		return []byte("New player joined")
	}
	return []byte("test function: " + event)
}

func (g * Game) run() {
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
				log.Print(string(message))
				for client := range g.clients {
					select {
					case client.send <- updateGameState(string(message)):
					default:
						close(client.send)
						delete(g.clients, client)
					}
				}
		}
	}
}