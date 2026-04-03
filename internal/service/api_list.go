package service

import (
	"back/internal/model"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type APIListService struct {
	db *gorm.DB
}

func NewAPIListService(db *gorm.DB) *APIListService {
	return &APIListService{db: db}
}

type CreateAPIListRequest struct {
	Currency string `json:"currency" binding:"required"`
	BaseURL  string `json:"base_url"`
	Code     string `json:"code"`
	Token    string `json:"token" binding:"required"`
}

type UpdateAPIListRequest struct {
	Currency *string `json:"currency"`
	BaseURL  *string `json:"base_url"`
	Code     *string `json:"code"`
	Token    *string `json:"token"`
}

func (s *APIListService) Create(in CreateAPIListRequest) (*model.APIList, error) {
	cc := normalizeCurrency(in.Currency)
	if cc == "" {
		return nil, errors.New("currency is required")
	}
	if strings.TrimSpace(in.Token) == "" {
		return nil, errors.New("token is required")
	}
	row := &model.APIList{
		Currency: cc,
		BaseURL:  strings.TrimSpace(in.BaseURL),
		Code:     strings.TrimSpace(in.Code),
		Token:    strings.TrimSpace(in.Token),
	}
	if err := s.db.Create(row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

func (s *APIListService) List() ([]model.APIList, error) {
	var rows []model.APIList
	if err := s.db.Order("id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *APIListService) GetByID(id uint) (*model.APIList, error) {
	var row model.APIList
	if err := s.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *APIListService) Update(id uint, in UpdateAPIListRequest) (*model.APIList, error) {
	var row model.APIList
	if err := s.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	if in.Currency != nil {
		cc := normalizeCurrency(*in.Currency)
		if cc == "" {
			return nil, errors.New("currency cannot be empty")
		}
		row.Currency = cc
	}
	if in.BaseURL != nil {
		row.BaseURL = strings.TrimSpace(*in.BaseURL)
	}
	if in.Code != nil {
		row.Code = strings.TrimSpace(*in.Code)
	}
	if in.Token != nil {
		if strings.TrimSpace(*in.Token) == "" {
			return nil, errors.New("token cannot be empty")
		}
		row.Token = strings.TrimSpace(*in.Token)
	}
	row.UpdatedAt = time.Now()
	if err := s.db.Save(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *APIListService) Delete(id uint) error {
	res := s.db.Delete(&model.APIList{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func normalizeCurrency(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
