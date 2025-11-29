package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/internal/server/service"
)

type Link struct {
	router fiber.Router

	link service.Link
}

func NewLink(router fiber.Router, service service.Service) *Link {
	api := &Link{
		router: router.Group("/link"),
		link:   *service.NewLink(),
	}

	api.createRoutes()

	return api
}

func (l *Link) createRoutes() {
	l.router.Get("/", l.getAll)
	l.router.Post("/sync", l.sync)
}

func (l *Link) getAll(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	links, err := l.link.GetAllByUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(links)
}

func (l *Link) sync(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var links []dto.Link

	if err := c.BodyParser(&links); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Var(links, "dive"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newLinks, err := l.link.Sync(c.Context(), userID, links)
	if err != nil {
		return err
	}

	return c.JSON(newLinks)
}
