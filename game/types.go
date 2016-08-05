package game

// Represents a game player
type Player struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Represents a song in the database
type Song struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Start string `json:"start"`
	Genre string `json:"genre"`
}

// Represents a song that has been picked
// with four title options that are the same genre
type PickResult struct {
	Song    *Song    `json:"song"`
	Options []string `json:"options"`
}

func NewPlayer(email, firstName, lastName string) *Player {
	return &Player{email, firstName, lastName}
}

func NewSong(id, title, start, genre string) *Song {
	return &Song{id, title, start, genre}
}

func NewPickResult(song *Song, options []string) *PickResult {
	return &PickResult{song, options}
}
