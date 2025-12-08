package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/sqlc"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Playlist struct {
	repo Repository

	history History
}

func (r *Repository) NewPlaylist() *Playlist {
	return &Playlist{
		repo:    *r,
		history: *r.NewHistory(),
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

func (p *Playlist) GetBySpotify(ctx context.Context, spotifyID string) (*model.Playlist, error) {
	playlist, err := p.repo.queries(ctx).PlaylistGetBySpotify(ctx, spotifyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get playlist by spotify id %s | %w", spotifyID, err)
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
		SpotifyID:     playlist.SpotifyID,
		OwnerUid:      playlist.OwnerUID,
		Name:          playlist.Name,
		Description:   pgtype.Text{String: playlist.Description, Valid: playlist.Description != ""},
		Public:        playlist.Public,
		TrackAmount:   int32(playlist.TrackAmount),
		Collaborative: playlist.Collaborative,
		CoverID:       pgtype.Text{String: playlist.CoverID, Valid: playlist.CoverID != ""},
		CoverUrl:      pgtype.Text{String: playlist.CoverURL, Valid: playlist.CoverURL != ""},
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

func (p *Playlist) CreateUser(ctx context.Context, user *model.PlaylistUser) error {
	id, err := p.repo.queries(ctx).PlaylistUserCreate(ctx, sqlc.PlaylistUserCreateParams{
		UserID:     int32(user.UserID),
		PlaylistID: int32(user.PlaylistID),
	})
	if err != nil {
		return fmt.Errorf("create playlist user %+v | %w", *user, err)
	}

	user.ID = int(id)

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
		CoverUrl:      pgtype.Text{String: playlist.CoverURL, Valid: playlist.CoverURL != ""},
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

func (p *Playlist) DeleteUserByUserPlaylist(ctx context.Context, user model.PlaylistUser) error {
	if err := p.repo.queries(ctx).PlaylistUserDeleteByUserPlaylist(ctx, sqlc.PlaylistUserDeleteByUserPlaylistParams{
		UserID:     int32(user.UserID),
		PlaylistID: int32(user.PlaylistID),
	}); err != nil {
		return fmt.Errorf("delete playlist user %+v | %w", user, err)
	}

	return nil
}
