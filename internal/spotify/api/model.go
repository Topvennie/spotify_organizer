package api

import (
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
)

type Album struct {
	SpotifyID   string `json:"id"`
	Name        string `json:"name"`
	TrackAmount int    `json:"total_tracks"`
	Popularity  int    `json:"popularity"`
}

func (a Album) ToModel() model.Album {
	return model.Album{
		SpotifyID:   a.SpotifyID,
		Name:        a.Name,
		TrackAmount: a.TrackAmount,
		Popularity:  a.Popularity,
	}
}

type Artist struct {
	SpotifyID string `json:"id"`
	Name      string `json:"name"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Popularity int `json:"popularity"`
}

func (a Artist) ToModel() model.Artist {
	return model.Artist{
		SpotifyID:  a.SpotifyID,
		Name:       a.Name,
		Followers:  a.Followers.Total,
		Popularity: a.Popularity,
	}
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Context struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type History struct {
	Track    Track     `json:"track"`
	PlayedAt time.Time `json:"played_at"`
	Context  Context   `json:"context"`
}

func (h History) ToModel(user model.User) model.History {
	return model.History{
		UserID:   user.ID,
		PlayedAt: h.PlayedAt,
	}
}

type Playlist struct {
	SpotifyID string `json:"id"`
	Owner     struct {
		UID         string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Tracks      struct {
		Total int `json:"total"`
	} `json:"tracks"`
	Collaborative bool    `json:"collaborative"`
	Images        []Image `json:"images"`
}

func (p *Playlist) ToModel() model.Playlist {
	url := ""
	maxWidth := -1
	for _, image := range p.Images {
		if image.Width > maxWidth {
			url = image.URL
			maxWidth = image.Width
		}
	}

	return model.Playlist{
		SpotifyID:     p.SpotifyID,
		OwnerUID:      p.Owner.UID,
		Name:          p.Name,
		Description:   p.Description,
		Public:        p.Public,
		TrackAmount:   p.Tracks.Total,
		Collaborative: p.Collaborative,
		CoverURL:      url,
		Owner: model.User{
			UID:         p.Owner.UID,
			DisplayName: p.Owner.DisplayName,
		},
	}
}

type Show struct {
	SpotifyID     string `json:"id"`
	Name          string `json:"name"`
	EpisodeAmount int    `json:"total_episodes"`
}

func (s Show) ToModel() model.Show {
	return model.Show{
		SpotifyID:     s.SpotifyID,
		Name:          s.Name,
		EpisodeAmount: s.EpisodeAmount,
	}
}

type Track struct {
	SpotifyID  string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

func (t *Track) ToModel() model.Track {
	return model.Track{
		SpotifyID:  t.SpotifyID,
		Name:       t.Name,
		Popularity: t.Popularity,
	}
}
