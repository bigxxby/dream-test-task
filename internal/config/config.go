package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte
var AppPort string

type Config struct {
	AppPort string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string

	JwtSecret string
}

// SetConfig reads the configuration from a JSON file and returns a Config struct
func SetConfig() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file")
	}
	return nil
}

func GetCofig() (*Config, error) {
	// Read the environment variables into the Config struct
	config := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSL_MODE"),
		AppPort:    os.Getenv("APP_PORT"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
	}

	// Check if any essential config is missing
	requiredConfigs := map[string]string{
		"DB_HOST":     config.DBHost,
		"DB_PORT":     config.DBPort,
		"DB_USER":     config.DBUser,
		"DB_PASSWORD": config.DBPassword,
		"DB_NAME":     config.DBName,
		"DB_SSL_MODE": config.DBSSLMode,
		"APP_PORT":    config.AppPort,
		"JWT_SECRET":  config.JwtSecret,
	}
	for key, value := range requiredConfigs {
		if value == "" {
			return nil, fmt.Errorf("missing required configuration value for %s", key)
		}
	}
	JwtSecret = []byte(config.JwtSecret)
	AppPort = config.AppPort
	return config, nil
}
