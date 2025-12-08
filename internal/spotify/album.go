package spotify

import (
	"context"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotify/api"
	"github.com/topvennie/sortifyr/pkg/utils"
)

// albumSync will syncronize the user's saved albums
func (c *client) albumSync(ctx context.Context, user model.User) error {
	albumsDB, err := c.album.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	albumsSpotifyAPI, err := c.api.AlbumGetAll(ctx, user)
	if err != nil {
		return err
	}
	albumsSpotify := utils.SliceMap(albumsSpotifyAPI, func(a api.Album) model.Album { return a.ToModel() })

	// Find albums we need to create
	for i := range albumsSpotify {
		if _, ok := utils.SliceFind(albumsDB, func(a *model.Album) bool { return a.Equal(albumsSpotify[i]) }); ok {
			continue
		}

		// User doesn't have this album yet
		album, err := c.album.GetBySpotify(ctx, albumsSpotify[i].SpotifyID)
		if err != nil {
			return err
		}
		if album == nil {
			// We don't have the album in our database yet
			if err := c.album.Create(ctx, &albumsSpotify[i]); err != nil {
				return err
			}
		}

		if err := c.album.CreateUser(ctx, &model.AlbumUser{UserID: user.ID, AlbumID: album.ID}); err != nil {
			return err
		}

		albumsDB = append(albumsDB, album)
	}

	// Find albums we need to delete
	for i := range albumsDB {
		if _, ok := utils.SliceFind(albumsSpotify, func(a model.Album) bool { return a.Equal(*albumsDB[i]) }); !ok {
			// User no longer has this album saved
			if err := c.album.DeleteUserByUserAlbum(ctx, model.AlbumUser{UserID: user.ID, AlbumID: albumsDB[i].ID}); err != nil {
				return err
			}
		}
	}

	return nil
}

// albumUpdate updates local album instances to match the spotify data
func (c *client) albumUpdate(ctx context.Context, user model.User) error {
	albumsDB, err := c.album.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	albumsSpotifyAPI, err := c.api.AlbumGetAll(ctx, user)
	if err != nil {
		return err
	}
	albumsSpotify := utils.SliceMap(albumsSpotifyAPI, func(a api.Album) model.Album { return a.ToModel() })

	for i := range albumsSpotify {
		albumDB, ok := utils.SliceFind(albumsDB, func(a *model.Album) bool { return a.Equal(albumsSpotify[i]) })
		if !ok {
			// Album not found
			if err := c.album.Create(ctx, &albumsSpotify[i]); err != nil {
				return err
			}

			continue
		}

		if !(*albumDB).EqualEntry(albumsSpotify[i]) {
			if err := c.album.Update(ctx, albumsSpotify[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *client) albumCheck(ctx context.Context, album *model.Album) error {
	albumDB, err := c.album.GetBySpotify(ctx, album.SpotifyID)
	if err != nil {
		return err
	}

	if albumDB == nil {
		return c.album.Create(ctx, album)
	}

	album.ID = albumDB.ID

	if !albumDB.EqualEntry(*album) {
		return c.album.Update(ctx, *album)
	}

	return nil
}
