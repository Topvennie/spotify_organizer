package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/sqlc"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

type Playlist struct {
	repo Repository
}

func (r *Repository) NewPlaylist() *Playlist {
	return &Playlist{
		repo: *r,
	}
}

func (p *Playlist) Get(ctx context.Context, playlistID int) (*model.Playlist, error) {
	playlist, err := p.repo.queries(ctx).PlaylistGet(ctx, int32(playlistID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get playlist by id %d | %w", playlistID, err)
	}

	return model.PlaylistModel(playlist), nil
}

func (p *Playlist) GetByUserPopulated(ctx context.Context, userID int) ([]*model.Playlist, error) {
	playlists, err := p.repo.queries(ctx).PlaylistGetByUserWithOwner(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get playlists by user %d | %w", userID, err)
	}

	return utils.SliceMap(playlists, func(r sqlc.PlaylistGetByUserWithOwnerRow) *model.Playlist {
		return model.PlaylistModelPopulated(r.Playlist, r.User)
	}), nil
}

func (p *Playlist) Create(ctx context.Context, playlist *model.Playlist) error {
	id, err := p.repo.queries(ctx).PlaylistCreate(ctx, sqlc.PlaylistCreateParams{
		UserID:        int32(playlist.UserID),
		SpotifyID:     playlist.SpotifyID,
		OwnerUid:      playlist.OwnerUID,
		Name:          playlist.Name,
		Description:   pgtype.Text{String: playlist.Description, Valid: playlist.Description != ""},
		Public:        playlist.Public,
		TrackAmount:   int32(playlist.TrackAmount),
		Collaborative: playlist.Collaborative,
		CoverID:       pgtype.Text{String: playlist.CoverID, Valid: playlist.CoverID != ""},
	})
	if err != nil {
		return fmt.Errorf("create playlist %+v | %w", *playlist, err)
	}

	playlist.ID = int(id)

	return nil
}

func (p *Playlist) CreateTrack(ctx context.Context, track *model.PlaylistTrack) error {
	id, err := p.repo.queries(ctx).PlaylistTrackCreate(ctx, sqlc.PlaylistTrackCreateParams{
		PlaylistID: int32(track.PlaylistID),
		TrackID:    int32(track.TrackID),
	})
	if err != nil {
		return fmt.Errorf("create playlist track %+v | %w", *track, err)
	}

	track.ID = int(id)

	return nil
}

func (p *Playlist) Update(ctx context.Context, playlist model.Playlist) error {
	if err := p.repo.queries(ctx).PlaylistUpdateBySpotify(ctx, sqlc.PlaylistUpdateBySpotifyParams{
		SpotifyID:     playlist.SpotifyID,
		OwnerUid:      playlist.OwnerUID,
		Name:          playlist.Name,
		Description:   pgtype.Text{String: playlist.Description, Valid: playlist.Description != ""},
		Public:        playlist.Public,
		TrackAmount:   int32(playlist.TrackAmount),
		Collaborative: playlist.Collaborative,
		CoverID:       pgtype.Text{String: playlist.CoverID, Valid: playlist.CoverID != ""},
	}); err != nil {
		return fmt.Errorf("update playlist %+v | %w", playlist, err)
	}

	return nil
}

func (p *Playlist) Delete(ctx context.Context, id int) error {
	if err := p.repo.queries(ctx).PlaylistDelete(ctx, int32(id)); err != nil {
		return fmt.Errorf("delete playlist %d | %w", id, err)
	}

	return nil
}

func (p *Playlist) DeleteTrackByPlaylistTrack(ctx context.Context, track model.PlaylistTrack) error {
	if err := p.repo.queries(ctx).PlaylistTrackDeleteByPlaylistTrack(ctx, sqlc.PlaylistTrackDeleteByPlaylistTrackParams{
		PlaylistID: int32(track.PlaylistID),
		TrackID:    int32(track.TrackID),
	}); err != nil {
		return fmt.Errorf("delete playlist track %+v | %w", track, err)
	}

	return nil
}
