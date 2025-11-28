// Package spotify connects with the spotify API
package spotify

import (
	"context"
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

func (c *client) playlistSync(ctx context.Context, user model.User) error {
	playlistsDB, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	playlistsSpotify, err := c.playlistGetAll(ctx, user)
	if err != nil {
		return err
	}

	toCreate := make([]model.Playlist, 0)
	toUpdate := make([]model.Playlist, 0)
	toDelete := make([]model.Playlist, 0)

	// Find the playlists that need to be created or updated
	for i := range playlistsSpotify {
		playlistDB, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.Equal(playlistsSpotify[i].model) })
		if !ok {
			// Playlist doesn't exist yet
			// Create it
			toCreate = append(toCreate, playlistsSpotify[i].model)
			if err := c.playlistSaveCover(&playlistsSpotify[i].model, nil, playlistsSpotify[i].Images); err != nil {
				return err
			}

			continue
		}

		// Regardless if any of the other data changed, let's update the cover if we can
		if err := c.playlistSaveCover(&playlistsSpotify[i].model, *playlistDB, playlistsSpotify[i].Images); err != nil {
			return err
		}

		// Playlist already exist
		// But is it still completely the same?
		if !(*playlistDB).EqualEntry(playlistsSpotify[i].model) {
			// Not completely the same anymore
			// Update it
			toUpdate = append(toUpdate, playlistsSpotify[i].model)
		}
	}

	for i := range toCreate {
		if err := c.playlistUserCheck(ctx, toCreate[i].OwnerUID); err != nil {
			return err
		}
		if err := c.playlist.Create(ctx, &toCreate[i]); err != nil {
			return err
		}
	}

	for i := range toUpdate {
		if err := c.playlistUserCheck(ctx, toUpdate[i].OwnerUID); err != nil {
			return err
		}
		if err := c.playlist.Update(ctx, toUpdate[i]); err != nil {
			return err
		}
	}

	// New and updated entries are now in the database
	// Let's bring our local copy up to date
	playlistsDB, err = c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	// Find the playlists that need to be deleted
	for _, playlistDB := range playlistsDB {
		_, ok := utils.SliceFind(playlistsSpotify, func(p playlist) bool { return p.model.SpotifyID == playlistDB.SpotifyID })
		if !ok {
			// Playlist no longer exists in the user's account
			// So delete it
			toDelete = append(toDelete, *playlistDB)
		}
	}

	for i := range toDelete {
		if err := c.playlist.Delete(ctx, toDelete[i].ID); err != nil {
			return err
		}
	}

	return nil
}

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
			Tracks:        p.Tracks.Total,
			Collaborative: p.Collaborative,
			Owner: model.User{
				UID:         p.Owner.UID,
				DisplayName: p.Owner.DisplayName,
			},
		},
		Images: p.Images,
	}
}

type playListResponse struct {
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

func (c *client) playlistGet(ctx context.Context, user model.User, limit, offset int) (playListResponse, error) {
	var resp playListResponse

	if err := c.request(ctx, user, fmt.Sprintf("me/playlists?offset=%d&limit=%d", offset, limit), &resp); err != nil {
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
