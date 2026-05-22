package config

import (
	"fmt"
	"time"
	// "time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Logging     LoggingConfig
	// CORS        CORSConfig
	SwaggerHost string // Base URL for Swagger docs (e.g., "api.swift-rms.org")
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port int
	Host string
	Env  string
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	URL            string
	MaxConnections int
	MinConnections int
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

 

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level  string
	Format string
}

// CORSConfig holds CORS configuration.
// type CORSConfig struct {
// 	AllowedOrigins []string
// 	AllowedMethods []string
// 	AllowedHeaders []string
// }


func LoadConfig() (*Config, error) {
	v := viper.New()
	
	// 1. Configure file details
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	// 2. Read the config file safely
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed reading config file: %w", err)
		}
	}

	// 3. Enable System Env overlays
	v.AutomaticEnv()

	// 4. Register defaults
	setDefaults(v)

	// 5. Construct the final object manually
	cfg := &Config{
		Server: ServerConfig{
			Env:  v.GetString("APP_ENV"),
			Port: v.GetInt("APP_PORT"),
		},
		Database: DatabaseConfig{
			URL:  v.GetString("DATABASE_URL"),
		},
		
		JWT: JWTConfig{
			Secret:        v.GetString("JWT_SECRET"),
			AccessExpiry:  v.GetDuration("JWT_ACCESS_EXPIRY"),
			RefreshExpiry: v.GetDuration("JWT_REFRESH_EXPIRY"),
			Issuer:        v.GetString("JWT_ISSUER"),
		},
		
	}

	// 6. Validate structural requirements before releasing the object
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}


func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", 8080)

	// v.SetDefault("DATABASE_MAX_CONNECTIONS", 10)
	// v.SetDefault("DATABASE_MIN_CONNECTIONS", 1)

	// Add more defaults as needed
}

func validate(cfg *Config) error {
	if cfg.Server.Env == "" {
		return fmt.Errorf("server environment (APP_ENV) is required")
	}
	if cfg.Server.Port == 0 {
		return fmt.Errorf("server port (APP_PORT) must be greater than 0")
	}
	if cfg.Database.URL == "" {
		return fmt.Errorf("database URL (DATABASE_URL) is required")
	}

	// Validate JWT configuration
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT secret (JWT_SECRET) is required")
	}
	if cfg.JWT.Issuer == "" {
		return fmt.Errorf("JWT issuer (JWT_ISSUER) is required")
	}
	if cfg.JWT.AccessExpiry <= 0 {
		return fmt.Errorf("JWT access expiry (JWT_ACCESS_EXPIRY) must be greater than 0")
	}
	if cfg.JWT.RefreshExpiry <= 0 {
		return fmt.Errorf("JWT refresh expiry (JWT_REFRESH_EXPIRY) must be greater than 0")
	}
	if cfg.JWT.RefreshExpiry <= cfg.JWT.AccessExpiry {
		return fmt.Errorf("JWT refresh expiry must be greater than access expiry")
	}

	return nil
}