package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/database/repository"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/pkg/utils"
	"go.uber.org/zap"
)

type Directory struct {
	service Service

	directory repository.Directory
}

func (s *Service) NewDirectory() *Directory {
	return &Directory{
		service:   *s,
		directory: *s.repo.NewDirectory(),
	}
}

func (d *Directory) GetByUser(ctx context.Context, userID int) ([]dto.Directory, error) {
	directoryModels, err := d.directory.GetByUserPopulated(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if directoryModels == nil {
		return []dto.Directory{}, nil
	}

	roots := utils.SliceFilter(directoryModels, func(d *model.Directory) bool { return d.ParentID == 0 })
	directories := make([]dto.Directory, 0, len(roots))

	for _, root := range roots {
		directories = append(directories, dto.DirectoryDTO(root, directoryModels))
	}

	return directories, nil
}

// Sync brings the database up to date with the data received from the api
func (d *Directory) Sync(ctx context.Context, userID int, roots []dto.Directory) ([]dto.Directory, error) {
	// Not the best practice to simply delete and recreate everything
	if err := d.service.withRollback(ctx, func(ctx context.Context) error {
		if err := d.directory.DeleteByUser(ctx, userID); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		for _, root := range roots {
			if err := d.create(ctx, userID, 0, root); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return d.GetByUser(ctx, userID)
}

// create is an internal function to create a single directory
// It goes recursively through the children to create every directory
func (d *Directory) create(ctx context.Context, userID, parentID int, directorySave dto.Directory) error {
	directory := directorySave.ToModel(userID, parentID)

	if err := d.directory.Create(ctx, directory); err != nil {
		return err
	}

	for _, child := range directorySave.Children {
		if err := d.create(ctx, userID, directory.ID, child); err != nil {
			return err
		}
	}

	return nil
}
