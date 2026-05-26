package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	config "github.com/roticeh/ipinfo/pkg/config"
)

// CorsGuard: Smart CORS that operates according to the rules in the config file.
func CorsGuard() fiber.Handler {
	allowedPatterns := config.AppConfig.Security.Cors.AllowedOrigins

	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return MatchOrigin(origin, allowedPatterns)
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
		MaxAge:           3600,
	})
}

// PublicCors: A simple CORS configuration that allows all origins and is used for public APIs where no authentication is required. This is suitable for endpoints that are meant to be accessed by any client without restrictions.
func PublicCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	})
}

// MatchOrigin checks whether the incoming origin matches any pattern in the pattern list.
// Supported Patterns:
// 1. “*”                  -> Allow everyone (Public)
// 2. “https://roticeh.com” -> Exact match
// 3. “http://localhost:*”  -> Localhost and any port (3000, 8080, etc.)
// 4. “https://*.roticeh.com” -> Only subdomains (api.roticeh.com) - EXCLUDING the main domain
// 5. “https://**.roticeh.com” -> (SPECIAL) Main domain AND all subdomains
func MatchOrigin(origin string, patterns []string) bool {

	if origin == "" {
		return false
	}

	for _, pattern := range patterns {
		// Public
		if pattern == "*" {
			return true
		}

		// Exact match
		if origin == pattern {
			return true
		}

		// Localhost ve Port Wildcard (http://localhost:*)
		if strings.HasPrefix(pattern, "http://localhost:") && strings.HasSuffix(pattern, "*") {
			base := strings.TrimSuffix(pattern, "*") // "http://localhost:"
			if strings.HasPrefix(origin, base) {
				return true
			}
		}

		//  "**.roticeh.com" -> Main Domain + Subdomains
		if strings.Contains(pattern, "**.") {
			// Pattern: https://**.roticeh.com
			// Clean:   https://roticeh.com
			cleanPattern := strings.Replace(pattern, "**.", "", 1)

			// Main Domain? (https://roticeh.com)
			if origin == cleanPattern {
				return true
			}

			// Subdomain? (https://api.roticeh.com)
			// Separate the protocol: “https://” and “roticeh.com”
			parts := strings.SplitN(cleanPattern, "://", 2)
			if len(parts) != 2 {
				continue // Malformed pattern, skip it
			}
			scheme := parts[0] // https
			host := parts[1]   // roticeh.com

			// Does the incoming origin start with the schema and end with the host?
			// Example: origin “https://api.roticeh.com”, suffix “.roticeh.com”
			if strings.HasPrefix(origin, scheme+"://") && strings.HasSuffix(origin, "."+host) {
				return true
			}
		}

		// 5. "*.roticeh.com" -> Only Subdomains (https://api.roticeh.com) - EXCLUDING the main domain (https://roticeh.com)
		if strings.Contains(pattern, "*.") {
			// Pattern: https://*.roticeh.com
			// Prefix: https://
			// Suffix: .roticeh.com
			parts := strings.Split(pattern, "*")
			if len(parts) == 2 {
				if strings.HasPrefix(origin, parts[0]) && strings.HasSuffix(origin, parts[1]) {
					return true
				}
			}
		}
	}

	return false
}
