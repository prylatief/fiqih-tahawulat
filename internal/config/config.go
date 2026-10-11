package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL    string
	Addr           string
	MaxOpenConns   int
	MaxIdleConns   int
	AllowedOrigins string
	LogLevel       string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("KUTUB_DATABASE_URL", "postgres://postgres:@localhost:5432/fiqih_tahawulat?sslmode=disable"),
		Addr:           getEnv("KUTUB_ADDR", ":3000"),
		MaxOpenConns:   getInt("KUTUB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:   getInt("KUTUB_MAX_IDLE_CONNS", 5),
		AllowedOrigins: getEnv("KUTUB_ALLOWED_ORIGINS", "http://localhost:3000"),
		LogLevel:       getEnv("KUTUB_LOG_LEVEL", "1"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
