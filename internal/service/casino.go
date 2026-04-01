package service

import (
	"back/internal/config"
	"back/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
		Token:        s.cfg.CasinoAPIKey,
		AgentCode:    s.cfg.CasinoAgent,
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
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal provider payload: %w", err)
	}

	retries := s.cfg.ProviderTries
	if retries < 1 {
		retries = 3
	}

	client := &http.Client{Timeout: 8 * time.Second}
	var lastErr error

	for attempt := 1; attempt <= retries; attempt++ {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.CasinoAPIURL, bytes.NewReader(body))
		if reqErr != nil {
			return fmt.Errorf("build provider request: %w", reqErr)
		}
		req.Header.Set("Content-Type", "application/json")

		start := time.Now()
		log.Printf("casino provider request attempt=%d method=POST url=%s", attempt, s.cfg.CasinoAPIURL)
		resp, doErr := client.Do(req)
		if doErr != nil {
			lastErr = doErr
			log.Printf("casino provider error attempt=%d err=%v", attempt, doErr)
		} else {
			respBody, readErr := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if readErr != nil {
				lastErr = readErr
			} else if resp.StatusCode >= 400 {
				lastErr = fmt.Errorf("provider status=%d body=%s", resp.StatusCode, string(respBody))
			} else if unmarshalErr := json.Unmarshal(respBody, out); unmarshalErr != nil {
				lastErr = fmt.Errorf("decode provider response: %w", unmarshalErr)
			} else {
				log.Printf("casino provider success attempt=%d latency=%s", attempt, time.Since(start))
				return nil
			}
		}

		if attempt < retries {
			time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
		}
	}

	return fmt.Errorf("provider request failed after retries: %w", lastErr)
}
