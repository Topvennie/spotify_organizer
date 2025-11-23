package spotify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/topvennie/spotify_organizer/pkg/config"
	"github.com/topvennie/spotify_organizer/pkg/redis"
)

type client struct {
	clientID     string
	clientSecret string
}

var C *client

func Init() error {
	clientID := config.GetString("auth.spotify.client.id")
	clientSecret := config.GetString("auth.spotify.client.secret")

	if clientID == "" || clientSecret == "" {
		return errors.New("client id or client secret not set")
	}

	C = &client{
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	return nil
}

func (c *client) NewUser(ctx context.Context, uid string, accessToken, refreshToken string, expiresIn time.Duration) error {
	if _, err := redis.C.Set(ctx, accessKey(uid), accessToken, expiresIn).Result(); err != nil {
		return fmt.Errorf("set access token %w", err)
	}

	if _, err := redis.C.Set(ctx, refreshKey(uid), refreshToken, 0).Result(); err != nil {
		return fmt.Errorf("set refresh token %w", err)
	}

	return nil
}

func accessKey(uid string) string {
	return uid + ":spotify:access_token"
}

func refreshKey(uid string) string {
	return uid + ":spotify:refresh_token"
}
