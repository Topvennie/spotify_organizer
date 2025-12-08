package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/topvennie/sortifyr/internal/database/model"
)

func (c *Client) AlbumGet(ctx context.Context, user model.User, spotifyID string) (Album, error) {
	var resp Album

	if err := c.request(ctx, user, http.MethodGet, "albums/"+spotifyID, http.NoBody, &resp); err != nil {
		return Album{}, fmt.Errorf("get album %s | %w", spotifyID, err)
	}

	return resp, nil
}

type albumAllResponse struct {
	Total int     `json:"total"`
	Items []Album `json:"items"`
}

func (c *Client) AlbumGetAll(ctx context.Context, user model.User) ([]Album, error) {
	albums := make([]Album, 0)

	total := 51
	limit := 50
	offset := 0

	for offset+limit < total {
		var resp albumAllResponse

		if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("me/albums?offset=%d&limit=%d", offset, limit), http.NoBody, &resp); err != nil {
			return nil, fmt.Errorf("get albums with limit %d and offset %d | %w", limit, offset, err)
		}

		albums = append(albums, resp.Items...)
		total = resp.Total

		offset += limit
	}

	return albums, nil
}
