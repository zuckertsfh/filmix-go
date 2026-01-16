package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	Mode      string

	Database   DatabaseConfig
	TmdbApiKey string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
	SSLMode      string
}

func Load() Config {
	// load .env file if exists
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	cfg := Config{
		Port:       getEnv("PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		Mode:       getEnv("APP_MODE", "development"),
		TmdbApiKey: getEnv("TMDB_API_KEY", ""),

		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASS", "password"),
			Name:         getEnv("DB_NAME", "appdb"),
			MaxOpenConns: getEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnv("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  getEnv("DB_MAX_IDLE_TIME", "15m"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
		},
	}

	if cfg.JWTSecret == "" {
		panic("JWT_SECRET must be set")
	}

	return cfg
}

func getEnv[T any](key string, defaultValue T) T {
	if value, exists := os.LookupEnv(key); exists {
		var zero T
		switch any(defaultValue).(type) {
		case int:
			// try parse int
			if i, err := strconv.Atoi(value); err == nil {
				return any(i).(T)
			}
		case string:
			return any(value).(T)
		default:
			return defaultValue
		}
		return zero // fallback if parsing fails
	}
	return defaultValue
}
