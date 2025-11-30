// Package spotify connects with the spotify API
package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/image"
	"github.com/topvennie/spotify_organizer/pkg/storage"
	"github.com/topvennie/spotify_organizer/pkg/utils"
	"go.uber.org/zap"
)

type playlistAPI struct {
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
	Collaborative bool               `json:"collaborative"`
	Images        []playlistImageAPI `json:"images"`
}

type playlistImageAPI struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type playlist struct {
	model  model.Playlist
	Images []playlistImageAPI
}

func (p *playlistAPI) toModel(user model.User) playlist {
	return playlist{
		model: model.Playlist{
			UserID:        user.ID,
			SpotifyID:     p.SpotifyID,
			OwnerUID:      p.Owner.UID,
			Name:          p.Name,
			Description:   p.Description,
			Public:        p.Public,
			TrackAmount:   p.Tracks.Total,
			Collaborative: p.Collaborative,
			Owner: model.User{
				UID:         p.Owner.UID,
				DisplayName: p.Owner.DisplayName,
			},
		},
		Images: p.Images,
	}
}

type playlistResponse struct {
	Total int           `json:"total"`
	Items []playlistAPI `json:"items"`
}

func (c *client) playlistGetAll(ctx context.Context, user model.User) ([]playlist, error) {
	playlists := make([]playlist, 0)

	limit := 50
	offset := 0

	resp, err := c.playlistGet(ctx, user, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get playlist with limit %d and offset %d | %w", limit, offset, err)
	}
	playlists = append(playlists, utils.SliceMap(resp.Items, func(p playlistAPI) playlist { return p.toModel(user) })...)

	total := resp.Total

	for offset+limit < total {
		offset += limit

		resp, err := c.playlistGet(ctx, user, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("get playlist with limit %d and offset %d | %w", limit, offset, err)
		}
		playlists = append(playlists, utils.SliceMap(resp.Items, func(p playlistAPI) playlist { return p.toModel(user) })...)
	}

	return playlists, nil
}

func (c *client) playlistGet(ctx context.Context, user model.User, limit, offset int) (playlistResponse, error) {
	var resp playlistResponse

	if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("me/playlists?offset=%d&limit=%d", offset, limit), http.NoBody, &resp); err != nil {
		return resp, fmt.Errorf("get playlist %w", err)
	}

	return resp, nil
}

// playlistUserCheck creates the user if it doesn't exist yet
func (c *client) playlistUserCheck(ctx context.Context, userUID string) error {
	user, err := c.user.GetByUID(ctx, userUID)
	if err != nil {
		return err
	}
	if user != nil {
		return nil
	}

	user = &model.User{
		UID: userUID,
	}

	if err := c.user.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

// playlistSaveCover will save and update covers for playlists
func (c *client) playlistSaveCover(newPlaylist, oldPlaylist *model.Playlist, images []playlistImageAPI) error {
	zap.S().Infof("Getting image for %s", newPlaylist.Name)
	if len(images) == 0 {
		return nil
	}

	// Get the biggest image
	var imageAPI *playlistImageAPI
	maxWidth := -1
	for _, i := range images {
		if i.Width > maxWidth {
			imageAPI = &i
			maxWidth = i.Width
		}
	}
	if imageAPI == nil || imageAPI.URL == "" {
		// No new image found
		return nil
	}

	resp, err := http.Get(imageAPI.URL)
	if err != nil {
		return fmt.Errorf("get image data %+v | %w", *newPlaylist, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read image data %+v | %w", *newPlaylist, err)
	}

	webp, err := image.ToWebp(data)
	if err != nil {
		return err
	}

	if oldPlaylist != nil && oldPlaylist.CoverID != "" {
		if err := storage.S.Delete(oldPlaylist.CoverID); err != nil {
			zap.S().Error(err) // Just log it, it's fine
		}
	}

	coverID := uuid.NewString()
	if err := storage.S.Set(coverID, webp, 0); err != nil {
		return fmt.Errorf("add cover image to storage %+v | %w", *newPlaylist, err)
	}

	newPlaylist.CoverID = coverID

	return nil
}

type playlistTrackAPI struct {
	Track struct {
		SpotifyID  string `json:"id"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
	} `json:"track"`
}

func (p *playlistTrackAPI) toModel() model.Track {
	return model.Track{
		SpotifyID:  p.Track.SpotifyID,
		Name:       p.Track.Name,
		Popularity: p.Track.Popularity,
	}
}

type playlistTrackResponse struct {
	Total int                `json:"total"`
	Items []playlistTrackAPI `json:"items"`
}

func (c *client) playlistGetTracksAll(ctx context.Context, user model.User, playlist model.Playlist) ([]model.Track, error) {
	tracks := make([]model.Track, 0)

	limit := 50
	offset := 0

	resp, err := c.playlistGetTracks(ctx, user, playlist, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get playlist tracks %+v with limit %d and offset %d | %w", playlist, limit, offset, err)
	}
	tracks = append(tracks, utils.SliceMap(resp.Items, func(t playlistTrackAPI) model.Track { return t.toModel() })...)

	total := resp.Total

	for offset+limit < total {
		offset += limit

		resp, err := c.playlistGetTracks(ctx, user, playlist, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("get playlist tracks %+v with limit %d and offset %d | %w", playlist, limit, offset, err)
		}
		tracks = append(tracks, utils.SliceMap(resp.Items, func(t playlistTrackAPI) model.Track { return t.toModel() })...)
	}

	return tracks, nil
}

func (c *client) playlistGetTracks(ctx context.Context, user model.User, playlist model.Playlist, limit, offset int) (playlistTrackResponse, error) {
	var resp playlistTrackResponse

	if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("playlists/%s/tracks?offset=%d&limit=%d", playlist.SpotifyID, offset, limit), http.NoBody, &resp); err != nil {
		return resp, fmt.Errorf("get playlist tracks %w", err)
	}

	return resp, nil
}

func (c *client) playlistAddTracksAll(ctx context.Context, user model.User, playlist model.Playlist, tracks []model.Track) error {
	current := 0
	total := len(tracks)

	for current < total {
		end := current + 100
		if end > total {
			end = total
		}

		toAdd := tracks[current:end]
		if err := c.playlistAddTracks(ctx, user, playlist, toAdd); err != nil {
			return fmt.Errorf("add tracks %d-%d to playlist %+v | %w", current, end, playlist, err)
		}

		current = end
	}

	return nil
}

func (c *client) playlistAddTracks(ctx context.Context, user model.User, playlist model.Playlist, tracks []model.Track) error {
	payload := struct {
		URIS []string `json:"uris"`
	}{
		URIS: utils.SliceMap(tracks, func(t model.Track) string { return "spotify:track:" + t.SpotifyID }),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal tracks payload: %w", err)
	}

	body := bytes.NewReader(data)

	if err := c.request(ctx, user, http.MethodPost, fmt.Sprintf("playlists/%s/tracks", playlist.SpotifyID), body, noResp); err != nil {
		return err
	}

	return nil
}
