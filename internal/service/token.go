package service

import (
	"errors"
	"time"

	"back/internal/model"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT creates a signed token containing the user's ID and role.
func GenerateJWT(user *model.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  float64(user.ID), // JWT parses numbers as float64
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// ValidateJWT parses and verifies a token using the provided secret.
// It returns the claims map if the token is valid.
func ValidateJWT(tokenString, secret string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
