package config

import "os"

type Config struct {
    DB *DBConfig
}

type DBConfig struct {
    DBUrl string
}

func New() *Config {
    return &Config{
        DB: &DBConfig{
            DBUrl: os.Getenv("DATABASE_URL"),
        },
    }
}
