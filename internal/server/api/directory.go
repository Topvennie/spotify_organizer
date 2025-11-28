package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/internal/server/service"
)

type Directory struct {
	router fiber.Router

	directory service.Directory
}

func NewDirectory(router fiber.Router, service service.Service) *Directory {
	api := &Directory{
		router:    router.Group("/directory"),
		directory: *service.NewDirectory(),
	}

	api.createRoutes()

	return api
}

func (d *Directory) createRoutes() {
	d.router.Get("/", d.get)
	d.router.Post("/sync", d.sync)
}

func (d *Directory) get(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	directories, err := d.directory.GetByUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(directories)
}

func (d *Directory) sync(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var directories []dto.Directory

	if err := c.BodyParser(&directories); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Var(directories, "dive"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newDirectories, err := d.directory.Sync(c.Context(), userID, directories)
	if err != nil {
		return err
	}

	return c.JSON(newDirectories)
}
