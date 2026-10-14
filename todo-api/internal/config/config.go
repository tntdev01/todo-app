package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

func required(key string) string {

	value := os.Getenv(key)

	if value == "" {
		panic(key + " is required")
	}

	return value
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		AppPort: required("APP_PORT"),

		DBHost:     required("DB_HOST"),
		DBPort:     required("DB_PORT"),
		DBName:     required("DB_NAME"),
		DBUser:     required("DB_USER"),
		DBPassword: required("DB_PASSWORD"),
		DBSSLMode:  required("DB_SSLMODE"),
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}
