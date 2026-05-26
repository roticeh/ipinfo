package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig contains HTTP server, networking, and base URL settings.
type ServerConfig struct {
	// Name is the service identifier (e.g., "IpInfo").
	Name string `mapstructure:"name"`

	// Port determines the network port the server listens on (e.g., 9980).
	Port int `mapstructure:"port"`

	// BodyLimit sets the maximum allowed request body size in Megabytes (MB) to prevent DoS.
	BodyLimit int `mapstructure:"body_limit"`

	// Timeouts configures server-side timeouts to protect against slow-client attacks.
	Timeouts TimeoutConfig `mapstructure:"timeouts"`

	// SecretKey is the master key for internal encryption and signatures.
	// WARNING: This is a critical security value.
	SecretKey string `mapstructure:"secret_key" validate:"required"`

	// BaseURL is the public-facing URL for the API (e.g., "https://ipinfo.roticeh.com").
	BaseURL string `mapstructure:"base_url" validate:"required,url"`
}

// TimeoutConfig defines precise timing constraints for HTTP connections.
type TimeoutConfig struct {
	// Read is the maximum duration for reading the entire request, including the body.
	Read time.Duration `mapstructure:"read"`
	// Write is the maximum duration before timing out writes of the response.
	Write time.Duration `mapstructure:"write"`
	// Idle is the maximum amount of time to wait for the next request when keep-alives are enabled.
	Idle time.Duration `mapstructure:"idle"`
}

// DatabaseConfig holds connection parameters for the primary data store.
type DatabaseConfig struct {
	// Path is the file path to the MaxMind GeoIP2 database for IP geolocation.
	Path string `mapstructure:"path"`
	// ASNPath is the file path to the MaxMind ASN database for Autonomous System Number lookups.
	ASNPath string `mapstructure:"asn_path"`
}

// SecurityConfig aggregates security-related policies including CORS, Rate Limiting, and Captcha.
type SecurityConfig struct {
	// Cors defines Cross-Origin Resource Sharing policies.
	Cors CorsConfig `mapstructure:"cors"`
	// RateLimit defines traffic throttling rules.
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// CorsConfig defines the rules for Cross-Origin requests.
type CorsConfig struct {
	// AllowedOrigins is a list of patterns (e.g. "https://*.roticeh.com") allowed to access the API.
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedHeaders string   `yaml:"allowed_headers"`
	AllowedMethods string   `yaml:"allowed_methods"`
	// AllowCredentials indicates whether the request can include user credentials like cookies.
	AllowCredentials bool `yaml:"allow_credentials"`
	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	MaxAge int `yaml:"max_age"`
}

// RateLimitConfig defines specific throttling limits for different route groups.
type RateLimitConfig struct {
	// Max is the maximum number of requests allowed within the window.
	Max int `mapstructure:"max"`
	// Expiration is the time window for the rate limit (e.g., "1m").
	Expiration time.Duration `mapstructure:"expiration"`
}
