package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUrl        string
	RedisAddr    string
	RateLimitRPM int
}

func Load() Config {
	return Config{
		DBUrl:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/saas?sslmode=disable"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RateLimitRPM: getEnvInt("RATE_LIMIT_RPM", 60),
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}