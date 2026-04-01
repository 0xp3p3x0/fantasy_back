package service

import (
	"back/internal/model"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AgentService struct {
	db *gorm.DB
}

type CreateAgentRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
	CallbackURL string `json:"callback_url" binding:"required,url"`
	Status      string `json:"status"`
}

type UpdateAgentRequest struct {
	Username    *string `json:"username"`
	Password    *string `json:"password"`
	Code        *string `json:"code"`
	Currency    *string `json:"currency"`
	CallbackURL *string `json:"callback_url"`
	Status      *string `json:"status"`
}

func NewAgentService(db *gorm.DB) *AgentService {
	return &AgentService{db: db}
}

func (s *AgentService) CreateAgent(in CreateAgentRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	agent := &model.User{
		Username:    in.Username,
		Password:    string(hash),
		Code:        in.Code,
		Currency:    in.Currency,
		CallbackURL: in.CallbackURL,
		Role:        "agent",
		Balance:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if in.Status == "" {
		agent.Status = "active"
	} else {
		agent.Status = in.Status
	}

	if err := s.db.Create(agent).Error; err != nil {
		return nil, err
	}
	return agent, nil
}

func (s *AgentService) ListAgents() ([]model.User, error) {
	var agents []model.User
	if err := s.db.Where("role = ?", "agent").Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (s *AgentService) UpdateAgent(id uint, in UpdateAgentRequest) (*model.User, error) {
	var agent model.User
	if err := s.db.Where("id = ? AND role = ?", id, "agent").First(&agent).Error; err != nil {
		return nil, err
	}

	if in.Username != nil && *in.Username != "" {
		agent.Username = *in.Username
	}
	if in.Code != nil && *in.Code != "" {
		agent.Code = *in.Code
	}
	if in.Currency != nil && *in.Currency != "" {
		agent.Currency = *in.Currency
	}
	if in.CallbackURL != nil && *in.CallbackURL != "" {
		agent.CallbackURL = *in.CallbackURL
	}
	if in.Status != nil && *in.Status != "" {
		agent.Status = *in.Status
	}
	if in.Password != nil && *in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		agent.Password = string(hash)
	}

	agent.UpdatedAt = time.Now()
	if err := s.db.Save(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentService) DeleteAgent(id uint) error {
	result := s.db.Where("id = ? AND role = ?", id, "agent").Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("agent not found")
	}
	return nil
}
