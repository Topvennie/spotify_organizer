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

type Link struct {
	repo Repository
}

func (r *Repository) NewLink() *Link {
	return &Link{
		repo: *r,
	}
}

func (l *Link) GetAllByUser(ctx context.Context, userID int) ([]*model.Link, error) {
	links, err := l.repo.queries(ctx).LinkGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get links by user %d | %w", userID, err)
	}

	return utils.SliceMap(links, model.LinkModel), nil
}

func (l *Link) Create(ctx context.Context, link *model.Link) error {
	id, err := l.repo.queries(ctx).LinkCreate(ctx, sqlc.LinkCreateParams{
		SourceDirectoryID: pgtype.Int4{Int32: int32(link.SourceDirectoryID), Valid: link.SourceDirectoryID != 0},
		SourcePlaylistID:  pgtype.Int4{Int32: int32(link.SourcePlaylistID), Valid: link.SourcePlaylistID != 0},
		TargetDirectoryID: pgtype.Int4{Int32: int32(link.TargetDirectoryID), Valid: link.TargetDirectoryID != 0},
		TargetPlaylistID:  pgtype.Int4{Int32: int32(link.TargetPlaylistID), Valid: link.TargetPlaylistID != 0},
	})
	if err != nil {
		return fmt.Errorf("create link %+v | %w", *link, err)
	}

	link.ID = int(id)

	return nil
}

func (l *Link) Update(ctx context.Context, link model.Link) error {
	if err := l.repo.queries(ctx).LinkUpdate(ctx, sqlc.LinkUpdateParams{
		ID:                int32(link.ID),
		SourceDirectoryID: pgtype.Int4{Int32: int32(link.SourceDirectoryID), Valid: link.SourceDirectoryID != 0},
		SourcePlaylistID:  pgtype.Int4{Int32: int32(link.SourcePlaylistID), Valid: link.SourcePlaylistID != 0},
		TargetDirectoryID: pgtype.Int4{Int32: int32(link.TargetDirectoryID), Valid: link.TargetDirectoryID != 0},
		TargetPlaylistID:  pgtype.Int4{Int32: int32(link.TargetPlaylistID), Valid: link.TargetPlaylistID != 0},
	}); err != nil {
		return fmt.Errorf("update link %+v | %w", link, err)
	}

	return nil
}

func (l *Link) Delete(ctx context.Context, linkID int) error {
	if err := l.repo.queries(ctx).LinkDelete(ctx, int32(linkID)); err != nil {
		return fmt.Errorf("delete link %d | %w", linkID, err)
	}

	return nil
}
