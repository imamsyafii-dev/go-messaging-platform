package config

import "os"

type Config struct {
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	JWTSecret  string
}

func Load() Config {
	return Config{
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5434"),
		DBUser:     getEnv("DB_USER", "messaging"),
		DBPassword: getEnv("DB_PASSWORD", "messaging"),
		DBName:     getEnv("DB_NAME", "messaging"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6380"),
		JWTSecret:  getEnv("JWT_SECRET", "development-only-secret"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
