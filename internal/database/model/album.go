package model

import "github.com/topvennie/sortifyr/pkg/sqlc"

type Album struct {
	ID          int
	SpotifyID   string
	Name        string
	TrackAmount int
	Popularity  int
	CoverID     string
	CoverURL    string
}

func AlbumModel(a sqlc.Album) *Album {
	coverID := ""
	if a.CoverID.Valid {
		coverID = a.CoverID.String
	}
	coverURL := ""
	if a.CoverUrl.Valid {
		coverURL = a.CoverUrl.String
	}

	return &Album{
		ID:          int(a.ID),
		SpotifyID:   a.SpotifyID,
		Name:        a.Name,
		TrackAmount: int(a.TrackAmount),
		Popularity:  int(a.Popularity),
		CoverID:     coverID,
		CoverURL:    coverURL,
	}
}

func (a *Album) Equal(a2 Album) bool {
	return a.SpotifyID == a2.SpotifyID
}

func (a *Album) EqualEntry(a2 Album) bool {
	return a.Name == a2.Name && a.TrackAmount == a2.TrackAmount && a.Popularity == a2.Popularity
}

type AlbumUser struct {
	ID      int
	UserID  int
	AlbumID int
}

func AlbumUserModel(a sqlc.AlbumUser) *AlbumUser {
	return &AlbumUser{
		ID:      int(a.ID),
		UserID:  int(a.UserID),
		AlbumID: int(a.AlbumID),
	}
}
