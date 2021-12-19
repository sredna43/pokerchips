package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", "localhost:8081", "http service address, eg 127.0.0.1:8081")

func main() {
	flag.Parse()
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.0.0"})
	r.LoadHTMLFiles("C:\\Users\\sredn\\Documents\\Code\\websites\\pokerchips\\client\\index.html")

	game := newGame()
	go game.run()

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		serveWs(game, c.Writer, c.Request)
	})
	r.Run(*addr)
}
