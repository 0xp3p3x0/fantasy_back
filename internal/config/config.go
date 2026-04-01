package config

import (
	"os"
	"sync"

	"gorm.io/gorm"
)

type Config struct {
	Port          string
	ProviderTries int
	SecretKey     string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
}

var (
	cfg  *Config
	once sync.Once
	db   *gorm.DB
)

func Load() (*Config, error) {
	cfg := &Config{
		Port:          getEnv("PORT", "3000"),
		ProviderTries: 3,
		SecretKey:     os.Getenv("SECRET_KEY"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getDB() *gorm.DB {
	return db
}
func SetDB(database *gorm.DB) {
	db = database
}

func LoadConfig() *Config {
	once.Do(func() {
		cfg = &Config{
			DBHost:        getEnv("DB_HOST", "localhost"),
			DBPort:        getEnv("DB_PORT", "5432"),
			DBUser:        getEnv("DB_USER", "postgres"),
			DBPassword:    getEnv("DB_PASSWORD", ""),
			DBName:        getEnv("DB_NAME", "casino_db"),
			Port:          getEnv("PORT", "3000"),
			ProviderTries: 3,
			SecretKey:     os.Getenv("SECRET_KEY"),
		}
	})
	return cfg
}
