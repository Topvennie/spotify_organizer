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

type Link struct {
	service Service

	link repository.Link
}

func (s *Service) NewLink() *Link {
	return &Link{
		service: *s,
		link:    *s.repo.NewLink(),
	}
}

func (l *Link) GetAllByUser(ctx context.Context, userID int) ([]dto.Link, error) {
	links, err := l.link.GetAllByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if links == nil {
		return []dto.Link{}, nil
	}

	return utils.SliceMap(links, dto.LinkDTO), nil
}

func (l *Link) Sync(ctx context.Context, userID int, linksSave []dto.Link) ([]dto.Link, error) {
	linksNew := utils.SliceMap(linksSave, func(l dto.Link) *model.Link { return l.ToModel() })

	linksDB, err := l.link.GetAllByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	toCreate := make([]*model.Link, 0)
	toUpdate := make([]*model.Link, 0)
	toDelete := make([]*model.Link, 0)

	for _, linkNew := range linksNew {
		if linkNew.ID == 0 {
			toCreate = append(toCreate, linkNew)
			continue
		}

		linkDB, ok := utils.SliceFind(linksDB, func(l *model.Link) bool { return l.ID == linkNew.ID })
		if !ok {
			// User gave an invalid id, unlucky
			continue
		}

		if !(*linkDB).Equal(*linkNew) {
			toUpdate = append(toUpdate, linkNew)
		}
	}

	for _, linkDB := range linksDB {
		if _, ok := utils.SliceFind(linksNew, func(l *model.Link) bool { return l.ID == linkDB.ID }); !ok {
			toDelete = append(toDelete, linkDB)
		}
	}

	// Do the database operations
	if err := l.service.withRollback(ctx, func(ctx context.Context) error {
		for _, link := range toCreate {
			if err := l.link.Create(ctx, link); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		for _, link := range toUpdate {
			if err := l.link.Update(ctx, *link); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		for _, link := range toDelete {
			if err := l.link.Delete(ctx, link.ID); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return l.GetAllByUser(ctx, userID)
}
