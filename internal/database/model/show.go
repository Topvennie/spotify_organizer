package model

import "github.com/topvennie/sortifyr/pkg/sqlc"

type Show struct {
	ID            int
	SpotifyID     string
	Name          string
	EpisodeAmount int
	CoverID       string
	CoverURL      string
}

func ShowModel(s sqlc.Show) *Show {
	coverID := ""
	if s.CoverID.Valid {
		coverID = s.CoverID.String
	}
	coverURL := ""
	if s.CoverUrl.Valid {
		coverURL = s.CoverUrl.String
	}

	return &Show{
		ID:            int(s.ID),
		SpotifyID:     s.SpotifyID,
		Name:          s.Name,
		EpisodeAmount: int(s.EpisodeAmount),
		CoverID:       coverID,
		CoverURL:      coverURL,
	}
}

func (s *Show) Equal(s2 Show) bool {
	return s.SpotifyID == s2.SpotifyID
}

func (s *Show) EqualEntry(s2 Show) bool {
	return s.Name == s2.Name && s.EpisodeAmount == s2.EpisodeAmount
}

type ShowUser struct {
	ID     int
	UserID int
	ShowID int
}

func ShowUserModel(s sqlc.ShowUser) *ShowUser {
	return &ShowUser{
		ID:     int(s.ID),
		UserID: int(s.UserID),
		ShowID: int(s.ShowID),
	}
}
