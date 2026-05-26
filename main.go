package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"

	server "github.com/roticeh/ipinfo/pkg/server"

	config "github.com/roticeh/ipinfo/pkg/config"
	logger "github.com/roticeh/ipinfo/pkg/logger"

	handlers "github.com/roticeh/ipinfo/pkg/handlers"
)

func main() {

	envFile := ".env"

	// Load the .env file
	envErr := godotenv.Load(envFile)
	if envErr != nil {
		if errors.Is(envErr, os.ErrNotExist) {

			logger.LogWarn("[NOTICE] No %s file found. Relying on system environment variables.", envFile)
		} else {
			logger.LogFatal("[ENV_ERROR] %s detected but parsing failed: %v", envFile, envErr)
		}
	}

	config.LoadConfig()

	var geoiploadErr error
	handlers.GeoCityDB, geoiploadErr = geoip2.Open(config.AppConfig.Database.Path)
	if geoiploadErr != nil {
		logger.LogFatal("Critical: City database could not be loaded: %v", geoiploadErr)
	}
	defer handlers.GeoCityDB.Close()

	// 3. ASN Veritabanını Başlat ve Handlers'a Enjekte Et
	handlers.GeoASNDB, geoiploadErr = geoip2.Open(config.AppConfig.Database.ASNPath)
	if geoiploadErr != nil {
		logger.LogFatal("Critical: ASN database could not be loaded: %v", geoiploadErr)
	}
	defer handlers.GeoASNDB.Close()

	app := server.StartServer()

	setupGracefulShutdown(app)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	logger.LogServerStart(config.AppConfig.Server.Port, config.AppConfig.Server.BaseURL)

	if err := app.Listen(addr); err != nil {
		logger.LogFatal("Server stopped unexpectedly: %v", err)
	}
}

func setupGracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit

		logger.LogInfo("[NOTICE] OS Shutdown signal received. Waiting for workers to finish...")

		// Stop Fiber server gracefully, allowing in-flight requests to complete.
		if err := app.Shutdown(); err != nil {
			logger.LogError("[SHUTDOWN_ERROR] Fiber shutdown failed: %v", err)
		}

		logger.LogInfo("[NOTICE] All workers finished. Shutting down securely.")
		os.Exit(0)
	}()
}
