package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleAgent = "agent"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey; autoIncrement"`
	Username  string    `json:"username" gorm:"unique; not null"`
	Password  string    `json:"password" gorm:"not null"`
	Currency  string    `json:"currency" gorm:"not null;default:USD"`
	CallbackURL string  `json:"callback_url" gorm:"not null;default:''"`
	Balance   float64   `json:"balance" gorm:"default:0"`
	Status    string    `json:"status" gorm:"not null;default:active"`
	Token     string    `json:"token" gorm:"type:varchar(36);uniqueIndex;not null"`
	Code      string    `json:"code" gorm:"unique;not null"`
	Role      string    `json:"role" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.Token == "" {
		u.Token = uuid.NewString()
	}
	return nil
}
