package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Telegram TelegramConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Host string
}

type JWTConfig struct {
	Secret string
}

type TelegramConfig struct {
	BotToken string
}

type AppConfig struct {
	Environment string
}

func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	// Parse database configuration
	dbConfig := parseDatabaseConfig()

	cfg := &Config{
		Database: dbConfig,
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", getEnv("PORT", "8080")),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		Telegram: TelegramConfig{
			BotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
		},
	}

	return cfg, nil
}

// parseDatabaseConfig parses database configuration from DATABASE_URL or individual env vars
func parseDatabaseConfig() DatabaseConfig {
	// Check for DATABASE_URL first (Railway, Heroku, etc.)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if parsed, err := parseDatabaseURL(dbURL); err == nil {
			return parsed
		}
	}

	// Fallback to individual environment variables
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", getEnv("PGHOST", "localhost")),
		Port:     getEnv("DB_PORT", getEnv("PGPORT", "5432")),
		User:     getEnv("DB_USER", getEnv("PGUSER", "flowpay")),
		Password: getEnv("DB_PASSWORD", getEnv("PGPASSWORD", "flowpay_password")),
		DBName:   getEnv("DB_NAME", getEnv("PGDATABASE", "flowpay_db")),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// parseDatabaseURL parses a DATABASE_URL into DatabaseConfig
func parseDatabaseURL(dbURL string) (DatabaseConfig, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return DatabaseConfig{}, err
	}

	password, _ := u.User.Password()
	dbName := strings.TrimPrefix(u.Path, "/")

	// Extract SSL mode from query parameters
	sslMode := "disable"
	if u.Query().Get("sslmode") != "" {
		sslMode = u.Query().Get("sslmode")
	} else if u.Scheme == "postgres" || u.Scheme == "postgresql" {
		// Railway usually requires SSL
		sslMode = "require"
	}

	return DatabaseConfig{
		Host:     u.Hostname(),
		Port:     u.Port(),
		User:     u.User.Username(),
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
