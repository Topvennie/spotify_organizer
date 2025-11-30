package model

import (
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
)

type Playlist struct {
	ID            int
	UserID        int
	SpotifyID     string
	OwnerUID      string
	Name          string
	Description   string
	Public        bool
	TrackAmount   int
	Collaborative bool
	CoverID       string

	// Non db fields
	Owner User
}

func PlaylistModel(p sqlc.Playlist) *Playlist {
	description := ""
	if p.Description.Valid {
		description = p.Description.String
	}
	coverID := ""
	if p.CoverID.Valid {
		coverID = p.CoverID.String
	}

	return &Playlist{
		ID:            int(p.ID),
		SpotifyID:     p.SpotifyID,
		OwnerUID:      p.OwnerUid,
		Name:          p.Name,
		Description:   description,
		Public:        p.Public,
		TrackAmount:   int(p.TrackAmount),
		Collaborative: p.Collaborative,
		CoverID:       coverID,
	}
}

func PlaylistModelPopulated(p sqlc.Playlist, u sqlc.User) *Playlist {
	playlist := PlaylistModel(p)
	playlist.Owner = *UserModel(u)

	return playlist
}

func (p *Playlist) Equal(p2 Playlist) bool {
	return p.SpotifyID == p2.SpotifyID
}

func (p *Playlist) EqualEntry(p2 Playlist) bool {
	return p.OwnerUID == p2.OwnerUID && p.Name == p2.Name && p.Description == p2.Description && p.Public == p2.Public && p.TrackAmount == p2.TrackAmount && p.Collaborative == p2.Collaborative && p.CoverID == p2.CoverID
}

type PlaylistTrack struct {
	ID         int
	PlaylistID int
	TrackID    int
}

func PlaylistTrackModel(p sqlc.PlaylistTrack) *PlaylistTrack {
	return &PlaylistTrack{
		ID:         int(p.ID),
		PlaylistID: int(p.PlaylistID),
		TrackID:    int(p.TrackID),
	}
}
