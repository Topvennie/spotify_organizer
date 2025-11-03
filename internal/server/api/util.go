package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func storeInSession(ctx *fiber.Ctx, key string, value any) error {
	session, err := goth_fiber.SessionStore.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session %w", err)
	}

	session.Set(key, value)

	return session.Save()
}
