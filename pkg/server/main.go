package server

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	config "github.com/roticeh/ipinfo/pkg/config"
	gerrors "github.com/roticeh/ipinfo/pkg/errors"
	logger "github.com/roticeh/ipinfo/pkg/logger"
	middleware "github.com/roticeh/ipinfo/pkg/middleware"
	routes "github.com/roticeh/ipinfo/pkg/routes"
	utils "github.com/roticeh/ipinfo/pkg/utils"
)

const DEFAULT_GLOBAL_SERVER_TIMEOUT time.Duration = 30 * time.Second

func StartServer() *fiber.App {

	cfg := config.AppConfig

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               cfg.Server.Name,
		ServerHeader:          "IpInfo",

		// Security & Performance Timeouts
		ReadTimeout:  cfg.Server.Timeouts.Read,
		WriteTimeout: cfg.Server.Timeouts.Write,
		IdleTimeout:  cfg.Server.Timeouts.Idle,

		// DDoS Protection (Body Size Limit)
		BodyLimit: cfg.Server.BodyLimit * 1024 * 1024,

		// Proxy Settings (Required for Load Balancers / Nginx)
		EnableTrustedProxyCheck: true,
		// TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader: "X-Forwarded-For",

		ErrorHandler: globalErrorHandler,
	})

	// Global Middlewares
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(helmet.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(middleware.GlobalTimeout(DEFAULT_GLOBAL_SERVER_TIMEOUT))

	// Initialize all routes
	// API Versioning: /v1
	apiGroup := app.Group("/v1")
	routes.SetupRoutes(apiGroup)

	app.Use(func(c *fiber.Ctx) error {
		return utils.ResponseError(
			c,
			fiber.StatusNotFound,
			gerrors.Err404NotFound,
			"The requested resource or endpoint does not exist on this server.",
		)
	})

	return app

}

// globalErrorHandler handles all unhandled errors in a unified JSON format.
func globalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := gerrors.ErrInternalServerError
	debugInfo := err.Error()

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
		debugInfo = fiberErr.Message
	}

	if code >= 500 {
		logger.LogError("Server Error: %s | Path: %s", debugInfo, c.Path())
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return utils.ResponseError(
			c,
			fiber.StatusRequestTimeout,
			gerrors.ErrTimeoutExceeded,
			"The system is currently experiencing high load and could not process your request in time. Please try again later.",
		)
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"code":    message,
		"message": debugInfo,
	})
}
