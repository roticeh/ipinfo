package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/roticeh/ipinfo/pkg/handlers"
	middleware "github.com/roticeh/ipinfo/pkg/middleware"
)

func SetupRoutes(app fiber.Router) {
	ipGroup := app.Group("/ip", middleware.CorsGuard(), middleware.GeneralLimiter())

	// GET /v1/ip/me -> Analyzes the calling client's own IP and User-Agent data.
	ipGroup.Get("/me", handlers.GetMyIPInfo)


	// POST /v1/ip/bulk -> Resolves multiple IP addresses simultaneously within a single payload.
	ipGroup.Post("/bulk", handlers.GetBulkIPInfo)

	// GET /v1/ip/:ipaddress -> Analyzes a specific target IP address passed as a parameter.
	ipGroup.Get("/:ipaddress", handlers.GetSpecificIPInfo)

	// GET /v1/ip/:ipaddress/:field -> Extracts a specific data block (e.g., /location, /network, /device) dynamically.
	ipGroup.Get("/:ipaddress/:field", handlers.GetSpecificFieldInfo)

	
}