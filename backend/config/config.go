package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
	SMTP   SMTPConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type ServerConfig struct {
	Port string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func Load() (*Config, error) {
	godotenv.Load()

	accessExpiry, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "sensory_navigator"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "default-secret-key"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

