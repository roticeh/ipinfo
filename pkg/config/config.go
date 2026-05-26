package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	logger "github.com/roticeh/ipinfo/pkg/logger"
)

var AppConfig *Config

// LoadConfig reads configuration from infra/app.yaml and environment variables.
func LoadConfig() {
	v := viper.New()

	setDefaults(v)

	v.SetConfigName("config")

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("..")

	v.SetEnvPrefix("IPINFO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// set custom env bindings for critical secrets and overrides
	setupCustomEnvBindings(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.LogInfo(" config.yaml not found, relying on Environment Variables.")
		} else {
			logger.LogFatal("CRITICAL: Error reading config file: %s", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		logger.LogFatal("FATAL: Failed to parse configuration structure: %v", err)
	}

	AppConfig = &cfg
	logger.LogInfo("System configuration loaded successfully.")
}

func setupCustomEnvBindings(v *viper.Viper) {

	// Server & Security
	v.BindEnv("server.secret_key", "SECRET_TOKEN")

	// Database
	v.BindEnv("database.path", "GEO_DB_PATH")
	v.BindEnv("database.asn_path", "ASN_DB_PATH")
	// Port Override
	v.BindEnv("server.port", "APP_PORT")

}

func setDefaults(v *viper.Viper) {

	// Server
	v.SetDefault("server.port", 9755)
	v.SetDefault("server.body_limit", 5) // 10MB

	v.SetDefault("server.secret_key", "")

	// Server.Timeouts
	v.SetDefault("server.timeouts.read", 5*time.Second)
	v.SetDefault("server.timeouts.write", 10*time.Second)
	v.SetDefault("server.timeouts.idle", 60*time.Second)



	// Database
	v.SetDefault("database.path", "./db/GeoLite2-City.mmdb")
	v.SetDefault("database.asn_path", "./db/GeoLite2-ASN.mmdb")

	// Rate Limits
	v.SetDefault("security.rate_limit.max", 150)
	v.SetDefault("security.rate_limit.expiration", 1*time.Minute)
}
