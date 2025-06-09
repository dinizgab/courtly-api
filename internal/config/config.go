package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	API     *APIConfig
	DB      *DBConfig
	SMTP    *SMTPConfig
	OpenPix *OpenPixConfig
	Storage *StorageConfig
}

type APIConfig struct {
	Port      string
	JwtSecret []byte
}

type DBConfig struct {
	DBUrl string
}

type SMTPConfig struct {
	Email string
	Host  string
	Port  int
	User  string
	Pass  string
}

type OpenPixConfig struct {
	BaseURL string
	AppID   string
	Timeout int
}

type StorageConfig struct {
	ProjectURL string
	APIKey     string
}

func New() (*Config, error) {
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("config.New - invalid SMTP_PORT: %w", err)
	}

	return &Config{
		API: &APIConfig{
			Port:      os.Getenv("API_PORT"),
			JwtSecret: []byte(os.Getenv("JWT_SECRET")),
		},
		DB: &DBConfig{
			DBUrl: os.Getenv("DATABASE_URL"),
		},
		SMTP: &SMTPConfig{
			Email: os.Getenv("SMTP_EMAIL"),
			Host:  os.Getenv("SMTP_HOST"),
			Port:  smtpPort,
			User:  os.Getenv("SMTP_USER"),
			Pass:  os.Getenv("SMTP_PASS"),
		},
		OpenPix: &OpenPixConfig{
			BaseURL: os.Getenv("OPENPIX_BASE_URL"),
			AppID:   os.Getenv("OPENPIX_APP_ID"),
		},
		Storage: &StorageConfig{
			ProjectURL: os.Getenv("STORAGE_PROJECT_URL"),
			APIKey:     os.Getenv("STORAGE_API_KEY"),
		},
	}, nil
}
