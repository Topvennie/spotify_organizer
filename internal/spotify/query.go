package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/pkg/redis"
	"go.uber.org/zap"
)

const (
	apiAccount = "https://accounts.spotify.com/api/token"
	apiSpotify = "https://api.spotify.com/v1"
)

type accountResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (c *client) refreshToken(ctx context.Context, uid string) error {
	zap.S().Info("Refreshing spotify access token")

	refreshToken, err := redis.C.Get(ctx, refreshKey(uid)).Result()
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			return fmt.Errorf("get redis key %s | %w", refreshKey(uid), err)
		}
		return fmt.Errorf("user %s refresh token not found", uid)
	}

	basicAuth := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "Basic %s:%s", c.clientID, c.clientSecret))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiAccount, http.NoBody)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", basicAuth)

	req.Form.Set("grant_type", "refresh_token")
	req.Form.Set("refresh_token", refreshToken)

	form := &fiber.Args{}
	form.Add("grant_type", "client_credentials")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %s", resp.Status)
	}

	var account accountResponse
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return fmt.Errorf("decode account response %w", err)
	}

	if account.TokenType != "Bearer" {
		return fmt.Errorf("invalid token type %+v", account)
	}

	if _, err := redis.C.Set(ctx, accessKey(uid), account.AccessToken, time.Duration(account.ExpiresIn)*time.Second).Result(); err != nil {
		return fmt.Errorf("set access token %w", err)
	}

	if account.RefreshToken != "" {
		if _, err := redis.C.Set(ctx, refreshKey(uid), account.RefreshToken, 0).Result(); err != nil {
			return fmt.Errorf("set refresh token %w", err)
		}
	}

	return nil
}

func (c *client) request(ctx context.Context, uid, url string, target any) error {
	zap.S().Infof("do request for url %s", url)

	accessToken, err := redis.C.Get(ctx, accessKey(uid)).Result()
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			return fmt.Errorf("get redis key %s | %w", accessKey(uid), err)
		}

		if err := c.refreshToken(ctx, uid); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", apiSpotify, url), http.NoBody)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case 401:
		return errors.New("bad or expired token")

	case 403:
		return errors.New("bad oauth request")

	case 429:
		zap.S().Info("rate limit hit")
		time.Sleep(5 * time.Second)

		return c.request(ctx, uid, url, target)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode body to json %w", err)
	}

	return nil
}
