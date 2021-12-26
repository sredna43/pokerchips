package main

import (
	"flag"
	"math/rand"
	"time"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sredna43/pokerchips/models"
)

var addr = flag.String("addr", ":8081", "http service address, eg 127.0.0.1:8081")

func generateLobby(n int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r1.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.0.0"})
	r.LoadHTMLFiles("C:\\Users\\sredn\\Documents\\Code\\websites\\pokerchips\\client\\index.html")

	game := newGame()
	go game.run()

	r.GET("/new_game", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		lobbyId := generateLobby(3)
		Lobbies[lobbyId] = &models.Lobby{
			GameState: &models.GameState{
				Players: make(map[string]*models.Player),
			},
			Settings: &models.Settings{
				MaxPlayers: 4,
				InitialChips: 100,
			},
			CreatedOn: time.Now().Unix(),
		}
		for lobby := range Lobbies {
			if time.Now().Unix() - Lobbies[lobby].CreatedOn > 86400 {
				delete(Lobbies, lobby)
			}
		}
		c.JSON(200, lobbyId)
	})

	r.GET("/:lobbyId", func(c *gin.Context) {
		log.Printf("%#v", Lobbies)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		lobbyId := strings.ToUpper(c.Param("lobbyId"))
		if _, ok := Lobbies[lobbyId]; ok {
			c.JSON(200, "OK")
		} else {
			c.JSON(404, "Table " + strings.ToUpper(lobbyId) + " not found.")
		}
	})

	r.GET("/ws", func(c *gin.Context) {
		serveWs(game, c.Writer, c.Request)
	})
	r.Run(*addr)
}
