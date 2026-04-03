package config

import (
	"os"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type Config struct {
	Port            string
	ProviderTries   int
	ProviderAPIURL  string // optional fallback when APIList.base_url is empty (same endpoint for all currencies)
	SecretKey       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

var (
	cfg  *Config
	once sync.Once
	db   *gorm.DB
)

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "3000"),
		ProviderTries:  3,
		ProviderAPIURL: os.Getenv("PROVIDER_API_URL"),
		SecretKey:      os.Getenv("SECRET_KEY"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
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

func getEnvAsInt(key string, fallback int) int {
	s := os.Getenv(key)
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
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
			DBHost:         getEnv("DB_HOST", "localhost"),
			DBPort:         getEnv("DB_PORT", "5432"),
			DBUser:         getEnv("DB_USER", "postgres"),
			DBPassword:     getEnv("DB_PASSWORD", ""),
			DBName:         getEnv("DB_NAME", "casino_db"),
			Port:           getEnv("PORT", "3000"),
			ProviderTries:  3,
			ProviderAPIURL: os.Getenv("PROVIDER_API_URL"),
			SecretKey:      os.Getenv("SECRET_KEY"),
		}
	})
	return cfg
}
