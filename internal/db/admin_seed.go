package db

import (
	"back/internal/model"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminCode     = "admin"
	defaultAdminPassword = "123!@#"
)

// EnsureDefaultAdmin creates the built-in admin user if no user with role admin exists.
func EnsureDefaultAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Where("role = ?", model.RoleAdmin).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := &model.User{
		Username:    defaultAdminUsername,
		Password:    string(hash),
		Code:        defaultAdminCode,
		Role:        model.RoleAdmin,
		Currency:    "USD",
		CallbackURL: "",
		Balance:     0,
		Status:      "active",
		Token:       uuid.NewString(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	log.Printf("Seeded default admin user (username=%s, code=%s)", defaultAdminUsername, defaultAdminCode)
	return nil
}
