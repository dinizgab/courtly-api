package config

import "os"

type Config struct {
    API *APIConfig
    DB *DBConfig
}

type APIConfig struct {
    Port string
}

type DBConfig struct {
    DBUrl string
}

func New() *Config {
    return &Config{
        API: &APIConfig{
            Port: os.Getenv("API_PORT"),
        },
        DB: &DBConfig{
            DBUrl: os.Getenv("DATABASE_URL"),
        },
    }
}
