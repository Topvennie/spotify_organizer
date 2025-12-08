package spotify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/spotify/api"
	"github.com/topvennie/sortifyr/pkg/concurrent"
	"github.com/topvennie/sortifyr/pkg/storage"
	"github.com/topvennie/sortifyr/pkg/utils"
)

// showSync will syncronize the user's saved shows
func (c *client) showSync(ctx context.Context, user model.User) error {
	showsDB, err := c.show.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	showsSpotifyAPI, err := c.api.ShowGetAll(ctx, user)
	if err != nil {
		return err
	}
	showsSpotify := utils.SliceMap(showsSpotifyAPI, func(s api.Show) model.Show { return s.ToModel() })

	// Find shows we need to create
	for i := range showsSpotify {
		if _, ok := utils.SliceFind(showsDB, func(s *model.Show) bool { return s.Equal(showsSpotify[i]) }); ok {
			continue
		}

		// User doesn't have this show yet
		show, err := c.show.GetBySpotify(ctx, showsSpotify[i].SpotifyID)
		if err != nil {
			return err
		}
		if show == nil {
			// We don't have the show in our database yet
			if err := c.show.Create(ctx, &showsSpotify[i]); err != nil {
				return err
			}
		}

		if err := c.show.CreateUser(ctx, &model.ShowUser{UserID: user.ID, ShowID: show.ID}); err != nil {
			return err
		}

		showsDB = append(showsDB, show)
	}

	// Find shows we need to delete
	for i := range showsDB {
		if _, ok := utils.SliceFind(showsSpotify, func(s model.Show) bool { return s.Equal(*showsDB[i]) }); !ok {
			// User no longer has this show saved
			if err := c.show.DeleteUserByUserShow(ctx, model.ShowUser{UserID: user.ID, ShowID: showsDB[i].ID}); err != nil {
				return err
			}
		}
	}

	return nil
}

// showUpdate updates local show instances to match the spotify data
func (c *client) showUpdate(ctx context.Context, user model.User) error {
	showsDB, err := c.show.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	showsSpotifyAPI, err := c.api.ShowGetAll(ctx, user)
	if err != nil {
		return err
	}
	showsSpotify := utils.SliceMap(showsSpotifyAPI, func(s api.Show) model.Show { return s.ToModel() })

	for i := range showsSpotify {
		showDB, ok := utils.SliceFind(showsDB, func(s *model.Show) bool { return s.Equal(showsSpotify[i]) })
		if !ok {
			// Show not found
			if err := c.show.Create(ctx, &showsSpotify[i]); err != nil {
				return err
			}

			continue
		}

		if !(*showDB).EqualEntry(showsSpotify[i]) {
			if err := c.show.Update(ctx, showsSpotify[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *client) showCoverSync(ctx context.Context, user model.User) error {
	shows, err := c.show.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	wg := concurrent.NewLimitedWaitGroup(12)

	var mu sync.Mutex
	var errs []error

	for _, show := range shows {
		if show.CoverURL == "" {
			continue
		}

		wg.Go(func() {
			cover, err := c.api.ImageGet(ctx, show.CoverURL)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			if len(cover) == 0 {
				return
			}

			oldCover := []byte{}
			if show.CoverID != "" {
				oldCover, err = storage.S.Get(show.CoverID)
				if err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("get cover for %+v | %w", *show, err))
					mu.Unlock()
					return
				}
			}

			if bytes.Equal(cover, oldCover) {
				return
			}

			show.CoverID = uuid.NewString()
			if err := storage.S.Set(show.CoverID, cover, 0); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("store new cover %+v | %w", *show, err))
				mu.Unlock()
				return
			}

			if err := c.show.Update(ctx, *show); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
		})
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c *client) showCheck(ctx context.Context, show *model.Show) error {
	showDB, err := c.show.GetBySpotify(ctx, show.SpotifyID)
	if err != nil {
		return err
	}

	if showDB == nil {
		return c.show.Create(ctx, show)
	}

	show.ID = showDB.ID

	if !showDB.EqualEntry(*show) {
		return c.show.Update(ctx, *show)
	}

	return nil
}
