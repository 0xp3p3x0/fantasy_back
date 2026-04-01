package model

import "time"

type APIList struct {
	ID        uint      `json:"id" gorm:"primaryKey; autoIncrement"`
	Currency  string    `json:"currency" gorm:"not null" default:"USD"`
	Code      string    `json:"code" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
