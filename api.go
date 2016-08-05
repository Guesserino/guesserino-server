package main

import (
	"github.com/pusher/foobar-server/game"
	"github.com/pusher/pusher-http-go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Foobar Server"))
}

// Starts a new gamge (admin)
func gameHandler(w http.ResponseWriter, req *http.Request) {
	if gameManager.OngoingGame {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("There is already a game in progress"))
		return
	}

	// send a signal to start the game
	gameManager.StartChannel <- struct{}{}

	w.WriteHeader(http.StatusOK)
	return
}

// pusher presence channel auth
func pusherAuthHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Failed to read request body in pusher auth: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queryParams, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("Failed parse query params: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerEmail := queryParams.Get("email")
	if playerEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player := gameManager.FindPlayer(playerEmail)
	if player == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	presenceData := pusher.MemberData{
		UserId: player.Email,
		UserInfo: map[string]string{
			"first_name": player.FirstName,
			"last_name":  player.LastName,
		},
	}

	response, err := game.Pusher.AuthenticatePresenceChannel(body, presenceData)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Handle webhooks
func webhookHandler(w http.ResponseWriter, req *http.Request) {
	// case member_added
	//     IF a game is ongoing, add them to the queued players list
	//     ELSE add them to the current players list
	// case member_removed
	//     IF the player is in the queued list, remove them
	//     ELSE IF player is in current players list
	//         Remove them form current players list
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	webhook, err := game.Pusher.Webhook(req.Header, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := webhook.Events[0]
	// first event only
	switch event.Name {
	case "member_added":
		log.Println("Received new member added event")

		player := gameManager.FindPlayer(event.UserId)
		if player == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("New player requesting to be added to the game: %s", player.Email)

		gameManager.PlayerJoinChannel <- player

		w.WriteHeader(http.StatusOK)
		return
	case "member_removed":
		log.Println("Received new member removed event")

		player := gameManager.FindPlayer(event.UserId)
		if player == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gameManager.PlayerLeftChannel <- player

		w.WriteHeader(http.StatusOK)
		return
	}
}

func playerHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(gameManager.PlayerContents)
	return

}
