package spotify

import (
	"context"
	"fmt"
)

type Playlist struct {
	Name string `json:"name"`
}

type playListResponse struct {
	Total int        `json:"total"`
	Items []Playlist `json:"items"`
}

func (c *client) GetAllPlaylists(ctx context.Context, uid string) ([]Playlist, error) {
	playlists := make([]Playlist, 0)

	limit := 50
	offset := 0

	resp, err := c.getPlaylists(ctx, uid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get playlist with limit %d and offset %d | %w", limit, offset, err)
	}
	playlists = append(playlists, resp.Items...)

	total := resp.Total

	for offset+limit < total {
		offset += limit

		resp, err := c.getPlaylists(ctx, uid, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("get playlist with limit %d and offset %d | %w", limit, offset, err)
		}
		playlists = append(playlists, resp.Items...)
	}

	return playlists, nil
}

func (c *client) getPlaylists(ctx context.Context, uid string, limit, offset int) (playListResponse, error) {
	var resp playListResponse

	if err := c.request(ctx, uid, fmt.Sprintf("me/playlists?offset=%d&limit=%d", offset, limit), &resp); err != nil {
		return resp, fmt.Errorf("get playlist %w", err)
	}

	return resp, nil
}
