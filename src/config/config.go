package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	APIKEY     string
	WhatsAppURL string
	WhatsAppAPIKEY string
	WhatsAppInstance string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		APIKEY:     getEnv("API_KEY", ""),
		WhatsAppURL: getEnv("WHATSAPP_URL", ""),
		WhatsAppAPIKEY: getEnv("WHATSAPP_API_KEY", ""),
		WhatsAppInstance: getEnv("WHATSAPP_INSTANCE", ""),
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" {
		return nil, fmt.Errorf("DBUser or DBPassword not configured")
	}

	return cfg, nil
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
