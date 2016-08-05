package game

import (
	"github.com/pusher/pusher-http-go"
	"log"
	"os"
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
	DefaultGameTime = 15
)

// Manages games, including setting up new games
type Manager struct {
	// Channel to indicate a player has joined
	PlayerJoinChannel chan *Player
	// Channel to signal player has left
	PlayerLeftChannel chan *Player
	// Current store instance
	*Store
	// If there is currently a game in progress
	OngoingGame bool
	// Signal start of game
	StartChannel chan struct{}
	// Reference to the current game instance
	CurrentGame *Game
}

func NewManager() *Manager {
	pChan := make(chan *Player)
	lChan := make(chan *Player)
	startChannel := make(chan struct{})

	return &Manager{pChan, lChan, NewStore(), false, startChannel, nil}
}

func (m *Manager) Setup() {
	endChan := make(chan struct{})
	game := NewGame(endChan, m.Store)

	log.Print("Reading files..")
	m.ReadAndStorePlayers()
	m.ReadAndStoreSongs()

	m.CurrentGame = game

	// Create a new song picker
	songPicker := NewSongPicker(m.Songs)

	for {
		select {
		case player := <-m.PlayerJoinChannel:
			log.Printf("New player %s joined", player.Email)

			if m.OngoingGame {
				log.Print("Game in progress. Queuing player..")
				m.AddQueuedPlayer(player)

				game.triggerPusherEvent(GameOngoingEvent, "{}")
			} else {
				log.Println("Added player to current list")
				m.AddCurrentPlayer(player)
			}
		case player := <-m.PlayerLeftChannel:
			log.Printf("Player leaving game: %s", player.Email)

			if m.PlayerInQueuedList(player) {
				log.Println("Player found in queued list. Removing..")
				m.RemoveQueuedPlayer(player)
			} else {
				if m.PlayerInCurrentList(player) {
					log.Println("Player found in current list. Removing..")
					m.RemoveCurrentPlayer(player)
				}
			}
		case <-m.StartChannel:
			if len(m.CurrentPlayers) == 0 {
				log.Println("Cannot start a game without any players, aborting")
				continue
			}

			log.Println("STARTING NEW GAME")
			m.OngoingGame = true

			queuedPlayers := len(m.QueuedPlayers)
			// add queued players if there are any
			if queuedPlayers > 0 {
				log.Printf("Found %d queued players found", queuedPlayers)
				m.AddQueuedPlayersToCurrentPlayers()

				// clear queued players
				m.ClearQueuedPlayers()
			}

			songPick := songPicker.Pick()

			go game.Start(songPick)
		case <-endChan:
			log.Println("GAME ENDED")
			m.OngoingGame = false

			game.triggerPusherEvent(GameEndEvent, "{}")
		}
	}
}
