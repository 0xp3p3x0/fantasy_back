package db

import (
	"back/internal/config"
	"back/internal/model"
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// applyUserUniqueness runs after AutoMigrate. Drops nickname unique so nicknames can duplicate (display name only).
func applyUserUniqueness(db *gorm.DB) error {
	var exists bool
	if err := db.Raw(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')`).Scan(&exists).Error; err != nil || !exists {
		return err
	}
	_ = db.Exec(`DROP INDEX IF EXISTS idx_users_nickname`).Error
	return nil
}

func Init(cfg *config.Config) *gorm.DB {
	sslMode := "require"
	if cfg.DBHost == "localhost" || cfg.DBHost == "127.0.0.1" {
		sslMode = "disable"
	}

	// Safely encode username and password
	user := url.UserPassword(cfg.DBUser, cfg.DBPassword)

	// Build URL
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
		RawQuery: url.Values{
			"sslmode": []string{sslMode},
		}.Encode(),
	}

	dsn := dbURL.String()
	log.Printf("Connecting to DB: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	// Auto-migrate the tables we care about
	if err := db.AutoMigrate(&model.User{}, &model.APIList{}, &model.APIList{}); err != nil {
		log.Fatalf("failed to automigrate: %v", err)
	}

	return db
}
