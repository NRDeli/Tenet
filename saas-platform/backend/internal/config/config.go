package config

import "os"

type Config struct {
    DBUrl string
}

func Load() Config {
    return Config{
        DBUrl: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/saas?sslmode=disable"),
    }
}

func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}