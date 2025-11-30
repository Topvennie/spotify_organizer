package model

import "github.com/topvennie/spotify_organizer/pkg/sqlc"

type Track struct {
	ID         int
	SpotifyID  string
	Name       string
	Popularity int
}

func TrackModel(t sqlc.Track) *Track {
	return &Track{
		ID:         int(t.ID),
		SpotifyID:  t.SpotifyID,
		Name:       t.Name,
		Popularity: int(t.Popularity),
	}
}

func (t *Track) Equal(t2 Track) bool {
	return t.SpotifyID == t2.SpotifyID
}

func (t *Track) EqualEntry(t2 Track) bool {
	return t.Name == t2.Name && t.Popularity == t2.Popularity
}
