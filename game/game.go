package game

import (
	"encoding/json"
	"log"
	"time"
)

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
			// game has ended
			return
		}
	}
}

// Signal the start of a game via Pusher
func (g *Game) triggerPusherEvent(eventName string, data string) {
	log.Printf("Triggering %s with data %s on Pusher", eventName, data)
	Pusher.Trigger(GameChannel, eventName, data)
}
