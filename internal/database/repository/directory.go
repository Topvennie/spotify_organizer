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

type Directory struct {
	repo Repository
}

func (r *Repository) NewDirectory() *Directory {
	return &Directory{
		repo: *r,
	}
}

func (d *Directory) GetByUserPopulated(ctx context.Context, userID int) ([]*model.Directory, error) {
	directoriesDB, err := d.repo.queries(ctx).DirectoryGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("directory get by user populated %d | %w", userID, err)
	}

	directoryPlaylistsDB, err := d.repo.queries(ctx).DirectoryPlaylistGetByDirectory(ctx, utils.SliceMap(directoriesDB, func(d sqlc.Directory) int32 { return d.ID }))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get directory playlists by  directories %+v | %w", directoriesDB, err)
		}
		directoryPlaylistsDB = []sqlc.DirectoryPlaylist{}
	}

	directoryPlaylistMap := make(map[int][]int)
	for i := range directoryPlaylistsDB {
		playlists, ok := directoryPlaylistMap[int(directoryPlaylistsDB[i].DirectoryID)]
		if !ok {
			playlists = []int{}
		}

		playlists = append(playlists, int(directoryPlaylistsDB[i].PlaylistID))
		directoryPlaylistMap[int(directoryPlaylistsDB[i].DirectoryID)] = playlists
	}

	playlistsDB, err := d.repo.queries(ctx).PlaylistGetByUserWithOwner(ctx, int32(userID))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get playlists by user %d | %w", userID, err)
		}
		playlistsDB = []sqlc.PlaylistGetByUserWithOwnerRow{}
	}

	playlistMap := make(map[int]*model.Playlist)
	for i := range playlistsDB {
		playlist := model.PlaylistModelPopulated(playlistsDB[i].Playlist, playlistsDB[i].User)

		playlistMap[playlist.ID] = playlist
	}

	directoryMap := make(map[int]*model.Directory)
	for i := range directoriesDB {
		directory, ok := directoryMap[int(directoriesDB[i].ID)]
		if !ok {
			directory = model.DirectoryModel(directoriesDB[i])
		}

		if playlistsIDs, ok := directoryPlaylistMap[directory.ID]; ok {
			for _, playlistID := range playlistsIDs {
				if playlist, ok := playlistMap[playlistID]; ok {
					directory.Playlists = append(directory.Playlists, *playlist)
				}
			}
		}

		directoryMap[directory.ID] = directory
	}

	return utils.MapValues(directoryMap), nil
}

func (d *Directory) Create(ctx context.Context, directory *model.Directory) error {
	return d.repo.WithRollback(ctx, func(ctx context.Context) error {
		id, err := d.repo.queries(ctx).DirectoryCreate(ctx, sqlc.DirectoryCreateParams{
			UserID:   int32(directory.UserID),
			Name:     directory.Name,
			ParentID: pgtype.Int4{Int32: int32(directory.ParentID), Valid: directory.ParentID != 0},
		})
		if err != nil {
			return fmt.Errorf("create directory %+v | %w", *directory, err)
		}

		directory.ID = int(id)

		for i := range directory.Playlists {
			if err := d.CreatePlaylist(ctx, &model.DirectoryPlaylist{
				DirectoryID: directory.ID,
				PlaylistID:  directory.Playlists[i].ID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *Directory) CreatePlaylist(ctx context.Context, directory *model.DirectoryPlaylist) error {
	id, err := d.repo.queries(ctx).DirectoryPlaylistCreate(ctx, sqlc.DirectoryPlaylistCreateParams{
		DirectoryID: int32(directory.DirectoryID),
		PlaylistID:  int32(directory.PlaylistID),
	})
	if err != nil {
		return fmt.Errorf("create directory playlist %+v | %w", *directory, err)
	}

	directory.ID = int(id)

	return nil
}

func (d *Directory) Update(ctx context.Context, directory model.Directory) error {
	if err := d.repo.queries(ctx).DirectoryUpdate(ctx, sqlc.DirectoryUpdateParams{
		ID:       int32(directory.ID),
		Name:     directory.Name,
		ParentID: pgtype.Int4{Int32: int32(directory.ParentID), Valid: directory.ParentID != 0},
	}); err != nil {
		return fmt.Errorf("update directory %+v | %w", directory, err)
	}

	return nil
}

func (d *Directory) DeleteByUser(ctx context.Context, userID int) error {
	if err := d.repo.queries(ctx).DirectoryDeleteByUser(ctx, int32(userID)); err != nil {
		return fmt.Errorf("delete directories by user %d | %w", userID, err)
	}

	// Directory playlists are deleted by cascade

	return nil
}
