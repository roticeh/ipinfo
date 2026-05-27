package handlers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/roticeh/ipinfo/pkg/config"
)

var startTime = time.Now()

func HealthCheck(c *fiber.Ctx) error {
	// database conn check
	cityDBStatus := "disconnected"
	asnDBStatus := "disconnected"

	if GeoCityDB != nil {
		cityDBStatus = "connected"
	}
	if GeoASNDB != nil {
		asnDBStatus = "connected"
	}

	statusCode := fiber.StatusOK
	statusText := "healthy"

	if cityDBStatus == "disconnected" || asnDBStatus == "disconnected" {
		statusCode = fiber.StatusServiceUnavailable
		statusText = "degraded"
	}

	response := fiber.Map{
		"success": statusCode == fiber.StatusOK,
		"status":  statusText,
	}

	providedToken := c.Query("token")
	expectedToken := config.AppConfig.Server.SecretKey

	if providedToken != "" && providedToken != expectedToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: Invalid healthcheck token.",
		})
	}

	if providedToken == expectedToken && expectedToken != "" {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		response["uptime"] = time.Since(startTime).String()
		response["databases"] = fiber.Map{
			"city_db": cityDBStatus,
			"asn_db":  asnDBStatus,
		}
		response["system"] = fiber.Map{
			"goroutines":   runtime.NumGoroutine(),
			"memory_alloc": m.Alloc / 1024 / 1024,
			"sys_memory":   m.Sys / 1024 / 1024,
		}
		response["timestamp"] = time.Now().Unix()
	}

	return c.Status(statusCode).JSON(response)
}
