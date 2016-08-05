package game

import (
	"encoding/json"
	"github.com/pusher/pusher-http-go"
	"log"
	"os"
	"time"
)

var Pusher *pusher.Client

func init() {
	Pusher = &pusher.Client{
		AppId:  os.Getenv("PUSHER_APP_ID"),
		Key:    os.Getenv("PUSHER_KEY"),
		Secret: os.Getenv("PUSHER_SECRET"),
	}
}

const (
	GameChannel      = "presence-foobar-game"
	GameStartEvent   = "start"
	GameEndEvent     = "end"
	GameOngoingEvent = "ongoing"
)

// Represents a new game
// This is pretty much stateless
// except for the embedded `Store`
type Game struct {
	EndChannel chan struct{}
	*Store
}

func NewGame(endChan chan struct{}, store *Store) *Game {
	return &Game{endChan, store}
}

// start a new game
// with a passed in `songPick`
// This allows calling `Start` as many times without
// having to share state
func (g *Game) Start(songPick *PickResult) {
	defer func() {
		// signal end
		g.EndChannel <- struct{}{}
	}()

	songTrigger, err := json.Marshal(songPick)
	if err != nil {
		log.Printf("Failed to JSON encode song pick: %s", err.Error())
	}

	g.triggerPusherEvent(GameStartEvent, string(songTrigger))

	for {
		select {
		case <-time.After(DefaultGameTime * time.Second):
			st, err := json.Marshal(songPick.Song)
			if err != nil {
				log.Printf("Failed to JSON encode song result: %s", err.Error())
			}
			g.triggerPusherEvent(GameEndEvent, string(st))
			return
		}
	}
}

// Signal the start of a game via Pusher
func (g *Game) triggerPusherEvent(eventName string, data string) {
	log.Printf("Triggering %s with data %s on Pusher", eventName, data)

	_, err := Pusher.Trigger(GameChannel, eventName, data)
	if err != nil {
		log.Printf("Failed to trigger event: %s", err.Error())
	}
}
