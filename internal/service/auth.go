package service

import (
	"back/internal/config"
	"back/internal/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Currency    string `json:"currency"`
	CallbackURL string `json:"callback_url"`
	Role        string `json:"role" binding:"required" default:"agent"`
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

func (s *AuthService) Register(in RegisterRequest) (user *model.User, err error) {
	var u model.User
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hash)
	u.Username = in.Username
	u.Role = in.Role
	u.Code = in.Code
	u.CallbackURL = in.CallbackURL
	if in.Currency == "" {
		u.Currency = "USD"
	} else {
		u.Currency = in.Currency
	}
	if in.Role == model.RoleAgent {
		if in.Currency == "" {
			return nil, errors.New("currency is required for agent")
		}
		if in.CallbackURL == "" {
			return nil, errors.New("callback_url is required for agent")
		}
	}
	u.Balance = 0
	u.Status = "active"
	u.Token = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if err := s.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *AuthService) Login(username string, password string) (token string, user *model.User, err error) {
	var u model.User
	if err := s.db.Where("username = ?", username).First(&u).Error; err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	tok, err := GenerateJWT(&u, s.cfg.SecretKey)
	if err != nil {
		return "", nil, err
	}
	return tok, &u, nil
}
