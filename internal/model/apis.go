package model

import "time"

type APIList struct {
	ID        uint      `json:"id" gorm:"primaryKey; autoIncrement"`
	Currency  string    `json:"currency" gorm:"not null;uniqueIndex;default:USD"`
	BaseURL   string    `json:"base_url" gorm:"type:text"` // provider Launch Game / GetGameUrl HTTP endpoint for this currency
	Code      string    `json:"code" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
