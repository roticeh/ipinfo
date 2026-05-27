package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	config "github.com/roticeh/ipinfo/pkg/config"
	handlers "github.com/roticeh/ipinfo/pkg/handlers"
	middleware "github.com/roticeh/ipinfo/pkg/middleware"
)

func SetupRoutes(app fiber.Router) {
	ipGroup := app.Group("/ip")
	ipGroupWithGuard := ipGroup.Group("", middleware.CorsGuard(), middleware.GeneralLimiter())

	// GET /v1/ip/me -> Analyzes the calling client's own IP and User-Agent data.
	ipGroup.Get("/me", handlers.GetMyIPInfo, middleware.PublicCors(), middleware.GeneralLimiter(config.RateLimitConfig{
		Max:        25,
		Expiration: 1 * time.Minute,
	}))

	// POST /v1/ip/bulk -> Resolves multiple IP addresses simultaneously within a single payload.
	ipGroupWithGuard.Post("/bulk", handlers.GetBulkIPInfo)

	// GET /v1/ip/:ipaddress -> Analyzes a specific target IP address passed as a parameter.
	ipGroupWithGuard.Get("/:ipaddress", handlers.GetSpecificIPInfo)

	// GET /v1/ip/:ipaddress/:field -> Extracts a specific data block (e.g., /location, /network, /device) dynamically.
	ipGroupWithGuard.Get("/:ipaddress/:field", handlers.GetSpecificFieldInfo)

}
