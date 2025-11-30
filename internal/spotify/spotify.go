package spotify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/redis"
)

var ErrUnauthorized = errors.New("access and refresh token expired")

type client struct {
	directory repository.Directory
	link      repository.Link
	playlist  repository.Playlist
	setting   repository.Setting
	track     repository.Track
	user      repository.User

	clientID     string
	clientSecret string
}

var C *client

func Init(repo repository.Repository) error {
	clientID := config.GetString("auth.spotify.client.id")
	clientSecret := config.GetString("auth.spotify.client.secret")

	if clientID == "" || clientSecret == "" {
		return errors.New("client id or client secret not set")
	}

	C = &client{
		directory:    *repo.NewDirectory(),
		link:         *repo.NewLink(),
		playlist:     *repo.NewPlaylist(),
		setting:      *repo.NewSetting(),
		track:        *repo.NewTrack(),
		user:         *repo.NewUser(),
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	return nil
}

func (c *client) NewUser(ctx context.Context, user model.User, accessToken, refreshToken string, expiresIn time.Duration) error {
	if _, err := redis.C.Set(ctx, accessKey(user), accessToken, expiresIn).Result(); err != nil {
		return fmt.Errorf("set access token %w", err)
	}

	if _, err := redis.C.Set(ctx, refreshKey(user), refreshToken, 0).Result(); err != nil {
		return fmt.Errorf("set refresh token %w", err)
	}

	return nil
}

func (c *client) Sync(ctx context.Context, user model.User) error {
	if err := c.playlistSync(ctx, user); err != nil {
		return fmt.Errorf("sync playlists for user %+v | %w", user, err)
	}

	if err := c.playlistTrackSync(ctx, user); err != nil {
		return fmt.Errorf("sync playlist tracks for user %+v | %w", user, err)
	}

	if err := c.userSync(ctx, user); err != nil {
		return fmt.Errorf("sync users for user %+v | %w", user, err)
	}

	if err := c.linkSync(ctx, user); err != nil {
		return fmt.Errorf("sync links for user %+v | %w", user, err)
	}

	setting, err := c.setting.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}
	if setting == nil {
		return fmt.Errorf("no setting found for user %+v", user)
	}

	setting.LastUpdate = time.Now()

	if err := c.setting.Update(ctx, *setting); err != nil {
		return err
	}

	return nil
}

func accessKey(user model.User) string {
	return user.UID + ":spotify:access_token"
}

func refreshKey(user model.User) string {
	return user.UID + ":spotify:refresh_token"
}
