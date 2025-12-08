package spotify

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotify/api"
	"github.com/topvennie/sortifyr/pkg/storage"
	"github.com/topvennie/sortifyr/pkg/utils"
)

func (c *client) playlistSync(ctx context.Context, user model.User) error {
	playlistsDB, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	playlistsSpotifyAPI, err := c.api.PlaylistGetAll(ctx, user)
	if err != nil {
		return err
	}
	playlistsSpotify := utils.SliceMap(playlistsSpotifyAPI, func(p api.Playlist) model.Playlist { return p.ToModel() })

	// Find playlists we need to create
	for i := range playlistsSpotify {
		if _, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.Equal(playlistsSpotify[i]) }); ok {
			continue
		}

		// User doesn't have this playlist yet
		playlist, err := c.playlist.GetBySpotify(ctx, playlistsSpotify[i].SpotifyID)
		if err != nil {
			return err
		}
		if playlist == nil {
			// We don't have the playlist in our database yet
			if err := c.userCheck(ctx, playlistsSpotify[i].OwnerUID); err != nil {
				return err
			}
			if err := c.playlist.Create(ctx, &playlistsSpotify[i]); err != nil {
				return err
			}
			playlist = &playlistsSpotify[i]
		}

		if err := c.playlist.CreateUser(ctx, &model.PlaylistUser{UserID: user.ID, PlaylistID: playlist.ID}); err != nil {
			return err
		}

		playlistsDB = append(playlistsDB, playlist)

	}

	// Find playlists we need to delete
	for i := range playlistsDB {
		if _, ok := utils.SliceFind(playlistsSpotify, func(p model.Playlist) bool { return p.Equal(*playlistsDB[i]) }); !ok {
			// User no longer has this playlist saved
			if err := c.playlist.DeleteUserByUserPlaylist(ctx, model.PlaylistUser{UserID: user.ID, PlaylistID: playlistsDB[i].ID}); err != nil {
				return err
			}
		}
	}

	return nil
}

// playlistUpdate updates local playlist instances to match the spotify data
func (c *client) playlistUpdate(ctx context.Context, user model.User) error {
	playlistsDB, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return err
	}

	playlistsSpotifyAPI, err := c.api.PlaylistGetAll(ctx, user)
	if err != nil {
		return err
	}
	playlistsSpotify := utils.SliceMap(playlistsSpotifyAPI, func(p api.Playlist) model.Playlist { return p.ToModel() })

	for i := range playlistsSpotify {
		playlistDB, ok := utils.SliceFind(playlistsDB, func(p *model.Playlist) bool { return p.Equal(playlistsSpotify[i]) })
		if !ok {
			// Playlist not found
			if err := c.userCheck(ctx, playlistsSpotify[i].OwnerUID); err != nil {
				return err
			}
			if err := c.playlist.Create(ctx, &playlistsSpotify[i]); err != nil {
				return err
			}

			continue
		}

		if !(*playlistDB).EqualEntry(playlistsSpotify[i]) {
			if err := c.userCheck(ctx, playlistsSpotify[i].OwnerUID); err != nil {
				return err
			}
			if err := c.playlist.Update(ctx, playlistsSpotify[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *client) playlistCoverSync(ctx context.Context, user model.User) (string, error) {
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if playlists == nil {
		return "", nil
	}

	newCovers := 0

	for _, playlist := range playlists {
		if playlist.CoverURL == "" {
			continue
		}

		cover, err := c.getCover(*playlist)
		if err != nil {
			return "", err
		}
		if len(cover) == 0 {
			continue
		}

		oldCover := []byte{}
		if playlist.CoverID != "" {
			oldCover, err = storage.S.Get(playlist.CoverID)
			if err != nil {
				return "", fmt.Errorf("get cover for %+v | %w", *playlist, err)
			}
		}

		if bytes.Equal(cover, oldCover) {
			continue
		}

		playlist.CoverID = uuid.NewString()
		if err := storage.S.Set(playlist.CoverID, cover, 0); err != nil {
			return "", fmt.Errorf("store new cover %+v | %w", *playlist, err)
		}

		if err := c.playlist.Update(ctx, *playlist); err != nil {
			return "", err
		}

		newCovers++
	}

	return fmt.Sprintf("New Covers: %d", newCovers), nil
}

// playlistTrackSync brings the local database up to date with the songs for each playlist
func (c *client) playlistTrackSync(ctx context.Context, user model.User) (string, error) {
	playlists, err := c.playlist.GetByUserPopulated(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if playlists == nil {
		return "", nil
	}

	totalCreated := 0
	totalDeleted := 0

	for _, playlist := range playlists {
		tracksDB, err := c.track.GetByPlaylist(ctx, playlist.ID)
		if err != nil {
			return "", err
		}

		tracksSpotifyAPI, err := c.api.PlaylistGetTrackAll(ctx, user, playlist.SpotifyID)
		if err != nil {
			return "", err
		}
		tracksSpotify := utils.SliceMap(tracksSpotifyAPI, func(t api.Track) model.Track { return t.ToModel() })

		toCreate := make([]model.Track, 0)
		toDelete := make([]model.Track, 0)

		for _, trackSpotify := range tracksSpotify {
			if _, ok := utils.SliceFind(tracksDB, func(t *model.Track) bool { return t.Equal(trackSpotify) }); !ok {
				toCreate = append(toCreate, trackSpotify)
			}
		}

		for _, trackDB := range tracksDB {
			if _, ok := utils.SliceFind(tracksSpotify, func(t model.Track) bool { return t.Equal(*trackDB) }); !ok {
				toDelete = append(toDelete, *trackDB)
			}
		}

		// Do the db operations
		for _, track := range toCreate {
			if err := c.trackCheck(ctx, &track); err != nil {
				return "", err
			}

			if err := c.playlist.CreateTrack(ctx, &model.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackID:    track.ID,
			}); err != nil {
				return "", err
			}
		}

		for _, track := range toDelete {
			if err := c.playlist.DeleteTrackByPlaylistTrack(ctx, model.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackID:    track.ID,
			}); err != nil {
				return "", err
			}
		}

		totalCreated += len(toCreate)
		totalDeleted += len(toDelete)
	}

	return fmt.Sprintf("Created %d | Deleted %d", totalCreated, totalDeleted), nil
}

func (c *client) playlistCheck(ctx context.Context, playlist *model.Playlist) error {
	playlistDB, err := c.playlist.GetBySpotify(ctx, playlist.SpotifyID)
	if err != nil {
		return err
	}

	if playlistDB == nil {
		if err := c.userCheck(ctx, playlist.OwnerUID); err != nil {
			return err
		}
		return c.playlist.Create(ctx, playlist)
	}

	playlist.ID = playlistDB.ID

	if !playlistDB.EqualEntry(*playlist) {
		return c.playlist.Update(ctx, *playlist)
	}

	return nil
}
