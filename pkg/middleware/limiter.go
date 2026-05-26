package middleware

import (
	"github.com/roticeh/ipinfo/pkg/config"
	"github.com/roticeh/ipinfo/pkg/errors"
	"github.com/roticeh/ipinfo/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func GeneralLimiter(limiterconfig ...config.RateLimitConfig) fiber.Handler {
	cfg := config.AppConfig.Security.RateLimit

	if len(limiterconfig) >= 1 {
		cfg.Max = limiterconfig[0].Max
	}

	if len(limiterconfig) >= 1 {
		cfg.Expiration = limiterconfig[0].Expiration
	}

	return limiter.New(limiter.Config{
		Max:        cfg.Max,
		Expiration: cfg.Expiration,

		LimiterMiddleware: limiter.SlidingWindow{},

		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ResponseError(
				c,
				fiber.StatusTooManyRequests,
				errors.ErrTooManyAttempts,
				"Too many requests. Please wait a moment.",
			)
		},
	})
}
