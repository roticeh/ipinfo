package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GlobalTimeout enforces a strict time limit on all HTTP requests and database queries.
// If the database or an external API takes too long, the context is cancelled to prevent RAM exhaustion.
func GlobalTimeout(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(c.UserContext(), timeout)

		defer cancel()

		c.SetUserContext(ctx)

		return c.Next()
	}
}
