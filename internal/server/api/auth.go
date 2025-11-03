// Package api contains all api routes
package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/spotify"
	"github.com/shareed2k/goth_fiber"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/internal/server/service"
	"github.com/topvennie/spotify_organizer/pkg/config"
	"go.uber.org/zap"
)

type Auth struct {
	router fiber.Router
	user   service.User

	redirectURL string
}

func NewAuth(router fiber.Router, service service.Service) *Auth {
	goth.UseProviders(
		spotify.New(
			config.GetString("auth.spotify.client_id"),
			config.GetString("auth.spotify.client_secret"),
			config.GetString("auth.spotify.callback_url"),
			spotify.ScopePlaylistReadPrivate,
		),
	)

	api := &Auth{
		router:      router.Group("/auth"),
		user:        *service.NewUser(),
		redirectURL: config.GetDefaultString("auth.redirect_url", "/"),
	}

	api.routes()

	return api
}

func (r *Auth) routes() {
	r.router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	r.router.Get("/callback/:provider", r.loginCallback)
	r.router.Post("/logout", r.logout)
}

func (r *Auth) loginCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		zap.S().Errorf("Failed to complete user auth %v", err)
		return fiber.ErrInternalServerError
	}

	dtoUser, err := r.user.GetByUID(c.Context(), user.UserID)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			dtoUser = dto.User{
				Name:  user.Name,
				Email: user.Email,
				UID:   user.UserID,
			}
			dtoUser, err = r.user.Save(c.Context(), dtoUser)
		}
		if err != nil {
			return err
		}
	}

	if err := storeInSession(c, "userID", dtoUser.ID); err != nil {
		zap.S().Errorf("Failed to store user id in session %v", err)
		return fiber.ErrInternalServerError
	}
	if err = storeInSession(c, "spotifyID", dtoUser.UID); err != nil {
		zap.S().Errorf("Failed to store spotify id in session %v", err)
		return fiber.ErrInternalServerError
	}

	return c.Redirect(r.redirectURL)
}

func (r *Auth) logout(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		zap.S().Errorf("Failed to logout %v", err)
	}

	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("Failed to get session %v", err)
		return fiber.ErrInternalServerError
	}
	if err := session.Destroy(); err != nil {
		zap.S().Errorf("Failed to destroy %v", err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
