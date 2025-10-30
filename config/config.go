package config

import (
	"fmt"
	"os"
)

// AppConfig holds all application configurations.
type AppConfig struct {
	DB        *DBConfig
	JWTSecret string
	Admin     *AdminConfig
}

// AdminConfig holds the default admin user configuration.
type AdminConfig struct {
	Email    string
	Password string
	Nombre   string
}

// DBConfig represents the database configuration.
type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN returns the data source name string for connecting to the database.
func (c *DBConfig) DSN() string {
	if c.Driver == "sqlite" {
		return "/tmp/api_test.db" // Using a temporary file for SQLite
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// LoadConfig loads all configurations from environment variables.
func LoadConfig() *AppConfig {
	return &AppConfig{
		DB: &DBConfig{
			Driver:   getEnv("DB_DRIVER", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWTSecret: getEnv("JWT_SECRET", "a_secure_secret"),
		Admin: &AdminConfig{
			Email:    getEnv("ADMIN_EMAIL", "admin@example.com"),
			Password: getEnv("ADMIN_PASSWORD", "admin123"),
			Nombre:   getEnv("ADMIN_NAME", "Admin"),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
