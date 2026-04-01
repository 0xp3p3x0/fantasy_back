package service

import (
	"back/internal/config"
	"back/internal/model"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProfileService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewProfileService(db *gorm.DB, cfg *config.Config) *ProfileService {
	return &ProfileService{
		db:  db,
		cfg: cfg,
	}
}

func (s *ProfileService) GetProfileById(id string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ProfileService) GetProfileByCode(code string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("code = ?", code).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ProfileService) UpdateProfile(username string, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	user.Password = password
	return &user, nil
}

func (s *ProfileService) UpdateCallbackURLByID(id uint, callbackURL string) (*model.User, error) {
	if callbackURL == "" {
		return nil, errors.New("callback_url is required")
	}

	var user model.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	user.CallbackURL = callbackURL
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ProfileService) ChangePasswordByID(id uint, currentPassword string, newPassword string) error {
	if currentPassword == "" || newPassword == "" {
		return errors.New("current_password and new_password are required")
	}
	if len(newPassword) < 6 {
		return errors.New("new_password must be at least 6 characters")
	}

	var user model.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is invalid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	return s.db.Save(&user).Error
}

func (s *ProfileService) RegenerateToken(code string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("code = ?", code).First(&user).Error; err != nil {
		return nil, err
	}
	user.Token = uuid.NewString()
	return &user, nil
}
