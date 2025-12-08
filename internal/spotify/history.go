package spotify

import (
	"context"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotify/api"
)

func (c *client) historySync(ctx context.Context, user model.User) (string, error) {
	latest, err := c.history.GetLatest(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if latest == nil {
		latest = &model.History{}
	}

	historySpotify, err := c.api.PlayerGetHistory(ctx, user)
	if err != nil {
		return "", err
	}

	toCreate := make([]api.History, 0)

	for i := range historySpotify {
		if !latest.PlayedAt.Before(historySpotify[i].PlayedAt) {
			break
		}

		toCreate = append(toCreate, historySpotify[i])
	}

	// Create history
	for _, h := range toCreate {
		if err := c.historyOneSync(ctx, user, h); err != nil {
			return "", nil
		}
	}

	return fmt.Sprintf("Tracks added: %d", len(toCreate)), nil
}

func (c *client) historyOneSync(ctx context.Context, user model.User, history api.History) error {
	historyModel := history.ToModel(user)

	trackModel := history.Track.ToModel()
	if err := c.trackCheck(ctx, &trackModel); err != nil {
		return err
	}
	historyModel.TrackID = trackModel.ID

	contextSpotifyID := uriToID(history.Context.URI)

	switch history.Context.Type {
	case "album":
		album, err := c.api.AlbumGet(ctx, user, contextSpotifyID)
		if err != nil {
			return err
		}
		albumModel := album.ToModel()
		if err := c.albumCheck(ctx, &albumModel); err != nil {
			return err
		}
		historyModel.AlbumID = albumModel.ID
	case "artist":
		artist, err := c.api.ArtistGet(ctx, user, contextSpotifyID)
		if err != nil {
			return err
		}
		artistModel := artist.ToModel()
		if err := c.artistCheck(ctx, &artistModel); err != nil {
			return err
		}
		historyModel.ArtistID = artistModel.ID
	case "playlist":
		playlist, err := c.api.PlaylistGet(ctx, user, contextSpotifyID)
		if err != nil {
			return err
		}
		playlistModel := playlist.ToModel()
		if err := c.playlistCheck(ctx, &playlistModel); err != nil {
			return err
		}
		historyModel.PlaylistID = playlistModel.ID
	case "show":
		show, err := c.api.ShowGet(ctx, user, contextSpotifyID)
		if err != nil {
			return err
		}
		showModel := show.ToModel()
		if err := c.showCheck(ctx, &showModel); err != nil {
			return err
		}
		historyModel.ShowID = showModel.ID
	}

	if err := c.history.Create(ctx, &historyModel); err != nil {
		return err
	}

	return nil
}
