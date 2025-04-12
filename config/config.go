package config

import (
	"errors"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	JWTSecret          string
}

func (c *Config) GetJWTSecret() string {
	return c.JWTSecret
}

// get singletonn config
var configInstance *Config
var configOnce sync.Once
var configError error

func GetConfig() (*Config, error) {
	configOnce.Do(func() {
		configInstance, configError = LoadConfig()
	})
	if configError != nil {
		return nil, configError
	}
	return configInstance, nil
}

func LoadConfig() (*Config, error) {

	// Load config from environment variables
	config, err := loadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfigFromEnv() (*Config, error) {

	// load from .env file
	err := godotenv.Load("local.env")
	if err != nil {
		return nil, err
	}
	serverPort := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	if serverPort == "" || databaseURL == "" || googleClientID == "" || googleClientSecret == "" || googleRedirectURI == "" {
		return nil, errors.New("missing required environment variables")
	}

	return &Config{
		ServerPort:         serverPort,
		DatabaseURL:        databaseURL,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURI:  googleRedirectURI,
	}, nil
}

func (c *Config) Validate() error {
	if c.ServerPort == "" || c.DatabaseURL == "" || c.GoogleClientID == "" || c.GoogleClientSecret == "" || c.GoogleRedirectURI == "" {
		return errors.New("missing required config values")
	}
	return nil
}

func (c *Config) GetServerPort() string {
	return c.ServerPort
}

func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}

func (c *Config) GetGoogleClientID() string {
	return c.GoogleClientID
}

func (c *Config) GetGoogleClientSecret() string {
	return c.GoogleClientSecret
}

func (c *Config) GetGoogleRedirectURI() string {
	return c.GoogleRedirectURI
}
