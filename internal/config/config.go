package config

import "os"

type Config struct {
	API *APIConfig
	DB  *DBConfig
}

type APIConfig struct {
	Port      string
	JwtSecret []byte
}

type DBConfig struct {
	DBUrl string
}

func New() *Config {
	return &Config{
		API: &APIConfig{
			Port:      os.Getenv("API_PORT"),
			JwtSecret: []byte(os.Getenv("JWT_SECRET")),
		},
		DB: &DBConfig{
			DBUrl: os.Getenv("DATABASE_URL"),
		},
	}
}
