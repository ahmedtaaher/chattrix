package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	DB         DBConfig
	Server     ServerConfig
	Auth       AuthConfig
}

type DBConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	SSLMode    string
}

type ServerConfig struct {
	Port       int
}

type AuthConfig struct {
	JWTSecret  string
}

func LoadConfig() *Config {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}
	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: mustGetEnv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME", "chattrixdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Auth: AuthConfig{
			JWTSecret: mustGetEnv("JWT_SECRET"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
  if value := os.Getenv(key); value != "" {
    return value
  }
  return defaultValue
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
		log.Fatalf("Invalid integer value for %s: %s", key, valueStr)
	}
	return defaultValue
}