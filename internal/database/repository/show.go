package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/sqlc"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Show struct {
	repo Repository
}

func (r *Repository) NewShow() *Show {
	return &Show{
		repo: *r,
	}
}

func (s *Show) GetBySpotify(ctx context.Context, spotifyID string) (*model.Show, error) {
	show, err := s.repo.queries(ctx).ShowGetBySpotify(ctx, spotifyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get show by spotify %s | %w", spotifyID, err)
	}

	return model.ShowModel(show), nil
}

func (s *Show) GetByUser(ctx context.Context, userID int) ([]*model.Show, error) {
	shows, err := s.repo.queries(ctx).ShowGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get shows by user  %d | %w", userID, err)
	}

	return utils.SliceMap(shows, model.ShowModel), nil
}

func (s *Show) Create(ctx context.Context, show *model.Show) error {
	id, err := s.repo.queries(ctx).ShowCreate(ctx, sqlc.ShowCreateParams{
		SpotifyID:     show.SpotifyID,
		Name:          show.Name,
		EpisodeAmount: int32(show.EpisodeAmount),
	})
	if err != nil {
		return fmt.Errorf("create show %+v | %w", *show, err)
	}

	show.ID = int(id)

	return nil
}

func (s *Show) CreateUser(ctx context.Context, user *model.ShowUser) error {
	id, err := s.repo.queries(ctx).ShowUserCreate(ctx, sqlc.ShowUserCreateParams{
		UserID: int32(user.UserID),
		ShowID: int32(user.ShowID),
	})
	if err != nil {
		return fmt.Errorf("create show user %+v | %w", *user, err)
	}

	user.ID = int(id)

	return nil
}

func (s *Show) Update(ctx context.Context, show model.Show) error {
	if err := s.repo.queries(ctx).ShowUpdate(ctx, sqlc.ShowUpdateParams{
		ID:            int32(show.ID),
		Name:          show.Name,
		EpisodeAmount: int32(show.EpisodeAmount),
	}); err != nil {
		return fmt.Errorf("update show %+v | %w", show, err)
	}

	return nil
}

func (s *Show) DeleteUserByUserShow(ctx context.Context, user model.ShowUser) error {
	if err := s.repo.queries(ctx).ShowUserDeleteByUserShow(ctx, sqlc.ShowUserDeleteByUserShowParams{
		UserID: int32(user.UserID),
		ShowID: int32(user.ShowID),
	}); err != nil {
		return fmt.Errorf("delete show user %+v | %w", user, err)
	}

	return nil
}
