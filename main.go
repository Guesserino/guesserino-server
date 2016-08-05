package main

import (
	"github.com/pusher/foobar-server/game"
	"net/http"
)

var (
	gameManager *game.Manager
)

func init() {
	// one single game manager instance
	gameManager = game.NewManager()
	go gameManager.Setup()
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/game/start", gameHandler)
	http.HandleFunc("/game/players", playerHandler)

	// Pusher related endpoints
	http.HandleFunc("/pusher/auth", pusherAuthHandler)
	http.HandleFunc("/pusher/webhook", webhookHandler)

	http.ListenAndServe(":8080", nil)
}
