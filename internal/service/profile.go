package service

import (
	"back/internal/config"
	"back/internal/model"

	"github.com/google/uuid"
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

func (s *ProfileService) RegenerateToken(code string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("code = ?", code).First(&user).Error; err != nil {
		return nil, err
	}
	user.Token = uuid.NewString()
	return &user, nil
}
