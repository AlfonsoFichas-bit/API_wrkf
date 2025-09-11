
package config

import (
	"fmt"
	"os"
)

// AppConfig holds all application configurations.
type AppConfig struct {
	DB        *DBConfig
	JWTSecret string
}

// DBConfig represents the database configuration.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN returns the data source name string for connecting to the database.
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// LoadConfig loads all configurations from environment variables.
func LoadConfig() *AppConfig {
	return &AppConfig{
		DB: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWTSecret: getEnv("JWT_SECRET", "5a8e02f0b339f9b67f85aa6a5160b7e134e50246e77fc273d78de1af49cfe365"),
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
