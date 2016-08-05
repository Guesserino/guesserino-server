package game

const OptionSize = 4

// Picks a song and matches songs of the same genre
// to present a list of possible winners
// Pick songs using a round robin strategy
// At least we're guaranteed a unique song till len(songs)
type SongPicker struct {
	// Index of the last picked song
	LastPickedIndex int
	// List of songs
	Songs []*Song
}

func NewSongPicker(songs []*Song) *SongPicker {
	return &SongPicker{0, songs}
}

func (sp *SongPicker) Pick() *PickResult {
	pickedSong := sp.pickSong()
	return NewPickResult(pickedSong, pickedSong.Id, sp.findOptionsForChosenSong(pickedSong))
}

func (sp *SongPicker) pickSong() *Song {
	songIndexToPick := (sp.LastPickedIndex + 1) % len(sp.Songs)
	sp.LastPickedIndex = songIndexToPick

	// chosen song
	return sp.Songs[songIndexToPick]
}

func (sp *SongPicker) findOptionsForChosenSong(song *Song) []string {
	titleList := make([]string, 0)

	for _, s := range sp.Songs {
		if s.Genre == song.Genre {
			titleList = append(titleList, s.Title)
			if len(titleList) == OptionSize {
				break
			}
		}
	}

	return titleList
}
