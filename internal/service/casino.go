package service

import (
	"back/internal/config"
	"back/internal/model"
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type CasinoService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewCasinoService(db *gorm.DB, cfg *config.Config) *CasinoService {
	return &CasinoService{
		db:  db,
		cfg: cfg,
	}
}

func (s *CasinoService) GetGameURL(in *model.GetGameURLRequest) (*model.GetGameURLResponse, error) {
	if in == nil {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "request is required"}
	}
	if in.UserCode == "" || in.VendorCode == "" || in.CurrencyCode == "" {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "userCode, vendorCode, and currencyCode are required"}
	}

	payload := getGameURLProviderRequest{
		Method:       "GetGameUrl",
		UserCode:     in.UserCode,
		NickName:     in.Nickname,
		VendorCode:   in.VendorCode,
		GameCode:     in.GameCode,
		CurrencyCode: in.CurrencyCode,
		Language:     in.Language,
		Channel:      in.Channel,
		IsDemo:       in.IsDemo,
	}

	var out getGameURLProviderResponse
	if err := s.postWithRetry(context.Background(), payload, &out); err != nil {
		return nil, err
	}

	if out.Status != 0 {
		return nil, &AppError{
			StatusCode: http.StatusBadGateway,
			Message:    fmt.Sprintf("provider error: %s", out.Msg),
		}
	}

	return &model.GetGameURLResponse{
		LaunchURL: out.LaunchURL,
	}, nil
}

type getGameURLProviderRequest struct {
	Method       string `json:"method"`
	Token        string `json:"token"`
	AgentCode    string `json:"agentCode"`
	UserCode     string `json:"userCode"`
	NickName     string `json:"nickName,omitempty"`
	VendorCode   string `json:"vendorCode"`
	GameCode     string `json:"gameCode,omitempty"`
	CurrencyCode string `json:"currencyCode"`
	Language     string `json:"language,omitempty"`
	Channel      string `json:"channel,omitempty"`
	IsDemo       bool   `json:"isDemo,omitempty"`
}

type getGameURLProviderResponse struct {
	Status    int    `json:"status"`
	Msg       string `json:"msg"`
	LaunchURL string `json:"launchUrl"`
}

func (s *CasinoService) postWithRetry(ctx context.Context, payload interface{}, out interface{}) error {

	return fmt.Errorf("provider request failed after retries")
}
