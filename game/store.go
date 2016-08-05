package game

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	DbFilePath      = "./files/db.json"
	PlayersFilePath = "./files/players.json"
)

// This is very much an in memory store
// Restarting the server means loss of game related data
// However, songs and players are read from files
// on startup and persisted
//
// The notion of players isn't relevant as most of it
// is handled by the front end and via Pusher
// However, it is beneficial to keep track of them
// on the server side
//
// Like the game, there is only one instance of the store
type Store struct {
	// All players in the database
	AllPlayers []*Player
	// Current players in the game
	CurrentPlayers []*Player
	// Queued Players
	QueuedPlayers []*Player
	// Songs
	Songs []*Song
	// Player file contents
	PlayerContents []byte
}

// Initial store
func NewStore() *Store {
	songs := make([]*Song, 0)
	pc := make([]byte, 0)
	return &Store{make([]*Player, 0), make([]*Player, 0), make([]*Player, 0), songs, pc}
}

// Reset the queue
func (s *Store) ClearQueuedPlayers() {
	s.QueuedPlayers = nil
}

// Add player to queue
func (s *Store) AddQueuedPlayer(player *Player) {
	s.QueuedPlayers = append(s.QueuedPlayers, player)
}

// Move all players in the queue to the current game
func (s *Store) AddQueuedPlayersToCurrentPlayers() {
	for _, p := range s.QueuedPlayers {
		s.CurrentPlayers = append(s.CurrentPlayers, p)
	}
}

// Remove player from queue
func (s *Store) RemoveQueuedPlayer(player *Player) {
	s.removePlayer(s.QueuedPlayers, player)
}

// Generic remove player function
func (s *Store) removePlayer(from []*Player, player *Player) {
	for i, p := range from {
		if p.Email == player.Email {
			from = append(from[:i], from[i+1:]...)
			break
		}
	}
}

// Add a player to the current game
func (s *Store) AddCurrentPlayer(player *Player) {
	s.CurrentPlayers = append(s.CurrentPlayers, player)
}

// Remove player from the current game
func (s *Store) RemoveCurrentPlayer(player *Player) {
	s.removePlayer(s.CurrentPlayers, player)
}

// Generic function to see if player is any list
func (s *Store) playerInList(list []*Player, player *Player) bool {
	for _, p := range list {
		if p.Email == player.Email {
			return true
		}
	}

	return false
}

// Check if player is in the queue
func (s *Store) PlayerInQueuedList(player *Player) bool {
	return s.playerInList(s.QueuedPlayers, player)
}

// Check if player is in the current game
func (s *Store) PlayerInCurrentList(player *Player) bool {
	return s.playerInList(s.CurrentPlayers, player)
}

// Find a player from the player list
func (s *Store) FindPlayer(email string) *Player {
	for _, p := range s.AllPlayers {
		if p.Email == email {
			return p
		}
	}

	return nil
}

/////////////////////// FILES /////////////////////////

func (s *Store) ReadAndStoreSongs() {
	data, err := ioutil.ReadFile(DbFilePath)
	if err != nil {
		log.Panicf("Error when creating song list: %v", err.Error())
	}

	songs := make([]*Song, 0)
	json.Unmarshal(data, &songs)

	log.Printf("Parsed song file. Found %d songs.", len(songs))

	s.Songs = songs
}

func (s *Store) ReadAndStorePlayers() {
	data, err := ioutil.ReadFile(PlayersFilePath)
	if err != nil {
		log.Panicf("Error when creating player list: %v", err.Error())
	}

	// Store the read file bytes so we don't have to read it again
	// the /game/players endpoint returns the same JSON
	// so don't read the file again
	s.PlayerContents = data

	players := make([]*Player, 0)
	json.Unmarshal(data, &players)

	log.Printf("Parsed player file. Found %d players.", len(players))

	s.AllPlayers = players
}
