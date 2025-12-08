package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/topvennie/sortifyr/internal/database/model"
)

func (c *Client) ShowGet(ctx context.Context, user model.User, spotifyID string) (Show, error) {
	var resp Show

	if err := c.request(ctx, user, http.MethodGet, "shows/"+spotifyID, http.NoBody, &resp); err != nil {
		return Show{}, fmt.Errorf("get show %s | %w", spotifyID, err)
	}

	return resp, nil
}

type showAllResponse struct {
	Total int    `json:"total"`
	Items []Show `json:"items"`
}

func (c *Client) ShowGetAll(ctx context.Context, user model.User) ([]Show, error) {
	shows := make([]Show, 0)

	total := 51
	limit := 50
	offset := 0

	for offset+limit < total {
		var resp showAllResponse

		if err := c.request(ctx, user, http.MethodGet, fmt.Sprintf("me/shows?offset=%d&limit=%d", offset, limit), http.NoBody, &resp); err != nil {
			return nil, fmt.Errorf("get shows with limit %d and offset %d | %w", limit, offset, err)
		}

		shows = append(shows, resp.Items...)
		total = resp.Total

		offset += limit
	}

	return shows, nil
}
