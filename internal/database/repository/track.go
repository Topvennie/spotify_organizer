package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

type Track struct {
	repo Repository
}

func (r *Repository) NewTrack() *Track {
	return &Track{
		repo: *r,
	}
}

func (t *Track) GetBySpotify(ctx context.Context, spotifyID string) (*model.Track, error) {
	track, err := t.repo.queries(ctx).TrackGetBySpotify(ctx, spotifyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get track by spotify id %s | %w", spotifyID, err)
	}

	return model.TrackModel(track), nil
}

func (t *Track) GetByPlaylist(ctx context.Context, playlistID int) ([]*model.Track, error) {
	tracks, err := t.repo.queries(ctx).TrackGetByPlaylist(ctx, int32(playlistID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tracks by playlist %d | %w", playlistID, err)
	}

	return utils.SliceMap(tracks, model.TrackModel), nil
}

func (t *Track) Create(ctx context.Context, track *model.Track) error {
	id, err := t.repo.queries(ctx).TrackCreate(ctx, sqlc.TrackCreateParams{
		SpotifyID:  track.SpotifyID,
		Name:       track.Name,
		Popularity: int32(track.Popularity),
	})
	if err != nil {
		return fmt.Errorf("create track %+v | %w", *track, err)
	}

	track.ID = int(id)

	return nil
}

func (t *Track) UpdateBySpotify(ctx context.Context, track model.Track) error {
	if err := t.repo.queries(ctx).TrackUpdateBySpotify(ctx, sqlc.TrackUpdateBySpotifyParams{
		SpotifyID:  track.SpotifyID,
		Name:       track.Name,
		Popularity: int32(track.Popularity),
	}); err != nil {
		return fmt.Errorf("update track %+v | %w", track, err)
	}

	return nil
}
