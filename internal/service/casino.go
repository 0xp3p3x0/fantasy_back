package service

import (
	"back/internal/config"
	"back/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var casinoHTTPClient = &http.Client{Timeout: 45 * time.Second}

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

func (s *CasinoService) GetGameURL(ctx context.Context, in *model.GetGameURLRequest) (*model.GetGameURLResponse, error) {
	if in == nil {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "request is required"}
	}
	if in.UserCode == "" || in.VendorCode == "" || in.CurrencyCode == "" {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "userCode, vendorCode, and currencyCode are required"}
	}

	var agent model.User
	code := strings.TrimSpace(in.AgentCode)
	if code == "" {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "agentCode is required"}
	}
	if err := s.db.Where("code = ? AND role = ?", code, model.RoleAgent).First(&agent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{StatusCode: http.StatusNotFound, Message: "agent not found for this agentCode"}
		}
		return nil, err
	}
	if strings.TrimSpace(agent.Code) == "" {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "agent code is not set for this account"}
	}

	var apiRow model.APIList
	cc := strings.ToUpper(strings.TrimSpace(in.CurrencyCode))
	if err := s.db.Where("UPPER(TRIM(currency)) = ?", cc).First(&apiRow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("no provider API configured for currency %s", in.CurrencyCode)}
		}
		return nil, err
	}
	if strings.TrimSpace(apiRow.Token) == "" {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "provider token missing for this currency"}
	}

	endpoint := strings.TrimSpace(apiRow.BaseURL)
	if endpoint == "" {
		endpoint = strings.TrimSpace(s.cfg.ProviderAPIURL)
	}
	if endpoint == "" {
		return nil, &AppError{StatusCode: http.StatusBadRequest, Message: "provider API URL not configured (set api_lists.base_url or PROVIDER_API_URL)"}
	}

	payload := model.ProviderGetGameURLRequest{
		Method:         in.Method,
		Token:          apiRow.Token,
		AgentCode:      apiRow.Code,
		UserCode:       strings.Join([]string{in.UserCode, agent.Code}, "_"),
		Nickname:       in.Nickname,
		VendorCode:     in.VendorCode,
		GameCode:       in.GameCode,
		CurrencyCode:   in.CurrencyCode,
		Language:       in.Language,
		Channel:        in.Channel,
		FreeRounds:     in.FreeRounds,
		FreeRoundsCode: in.FreeRoundsCode,
		CustomGameName: in.CustomGameName,
		HomeUrl:        in.HomeUrl,
		IsDemo:         in.IsDemo,
		UserBalance:    in.UserBalance,
	}

	var out model.ProviderGetGameURLResponse
	if err := s.postWithRetry(ctx, endpoint, payload, &out); err != nil {
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

func (s *CasinoService) postWithRetry(ctx context.Context, endpoint string, payload interface{}, out interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	tries := s.cfg.ProviderTries
	if tries < 1 {
		tries = 1
	}

	var lastErr error
	for attempt := 0; attempt < tries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := casinoHTTPClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
			continue
		}

		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode provider response: %w", err)
		}
		return nil
	}

	if lastErr != nil {
		return &AppError{
			StatusCode: http.StatusBadGateway,
			Message:    fmt.Sprintf("provider request failed after retries: %v", lastErr),
		}
	}
	return &AppError{
		StatusCode: http.StatusBadGateway,
		Message:    "provider request failed after retries",
	}
}

// HandleProviderCallback dispatches one shared callback endpoint by req.Method.
func (s *CasinoService) HandleProviderCallback(ctx context.Context, req *model.CallbackRequest) *model.ProviderCallbackResponse {
	if req == nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INVALID_REQUEST"}
	}
	method := strings.TrimSpace(req.Method)
	switch {
	case strings.EqualFold(method, "GetBalance"):
		return s.callbackGetBalance(ctx, req)
	case strings.EqualFold(method, "ChangeBalance"):
		return s.callbackChangeBalance(ctx, req)
	case strings.EqualFold(method, "UpdateDetail"):
		return s.callbackUpdateDetail(ctx, req)
	default:
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "UNKNOWN_METHOD"}
	}
}

func (s *CasinoService) callbackGetBalance(ctx context.Context, req *model.CallbackRequest) *model.ProviderCallbackResponse {
	if strings.TrimSpace(req.Token) == "" || strings.TrimSpace(req.UserCode) == "" || strings.TrimSpace(req.CurrencyCode) == "" {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INVALID_REQUEST"}
	}
	if _, err := s.findAPIListByTokenAndCurrency(req.Token, req.CurrencyCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ProviderCallbackResponse{Status: model.ProvCallbackInvalidToken, Msg: "INVALID_TOKEN"}
		}
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INTERNAL_ERROR"}
	}
	_, player, err := s.resolvePlayerFromCallbackUserCode(req.UserCode)
	if err != nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackInvalidUser, Msg: "INVALID_USER"}
	}

	bal, err := s.getPlayerBalance(player.ID, strings.ToUpper(strings.TrimSpace(req.CurrencyCode)))
	if err != nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INTERNAL_ERROR"}
	}
	b := bal
	return &model.ProviderCallbackResponse{Status: model.ProvCallbackOK, Msg: "SUCCESS", Balance: &b}
}

func (s *CasinoService) callbackChangeBalance(ctx context.Context, req *model.CallbackRequest) *model.ProviderCallbackResponse {
	if err := reqChangeBalanceOK(req); err != nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INVALID_REQUEST"}
	}
	if _, err := s.findAPIListByTokenAndCurrency(req.Token, req.CurrencyCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ProviderCallbackResponse{Status: model.ProvCallbackInvalidToken, Msg: "INVALID_TOKEN"}
		}
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INTERNAL_ERROR"}
	}

	agent, player, err := s.resolvePlayerFromCallbackUserCode(req.UserCode)
	if err != nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackInvalidUser, Msg: "INVALID_USER"}
	}

	cc := strings.ToUpper(strings.TrimSpace(req.CurrencyCode))
	txnCode := strings.TrimSpace(req.TxnCode)
	wagerID := *req.WagerID
	amount := *req.Amount

	var outBal float64
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.CasinoProcessedTxn
		if err := tx.Where("txn_code = ?", txnCode).First(&existing).Error; err == nil {
			// Idempotent retry: return current balance
			b, e := s.lockedPlayerBalance(tx, player.ID, cc)
			if e != nil {
				return e
			}
			outBal = b.Balance
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		balRow, err := s.lockedPlayerBalance(tx, player.ID, cc)
		if err != nil {
			return err
		}
		newBal := balRow.Balance + amount
		if newBal < 0 {
			return errInsufficient
		}
		if err := tx.Model(&model.CasinoPlayerBalance{}).
			Where("player_id = ? AND currency_code = ?", player.ID, cc).
			Update("balance", newBal).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.CasinoProcessedTxn{
			TxnCode:      txnCode,
			WagerID:      wagerID,
			PlayerID:     player.ID,
			CurrencyCode: cc,
			BalanceAfter: newBal,
		}).Error; err != nil {
			return err
		}

		var w model.CasinoWager
		wErr := tx.Where("wager_id = ?", wagerID).First(&w).Error
		if errors.Is(wErr, gorm.ErrRecordNotFound) {
			w = model.CasinoWager{
				WagerID:      wagerID,
				PlayerID:     player.ID,
				AgentID:      agent.ID,
				CurrencyCode: cc,
				VendorCode:   strings.TrimSpace(req.VendorCode),
				GameCode:     strings.TrimSpace(req.GameCode),
				GameRoundID:  strings.TrimSpace(req.GameRoundID),
			}
			if err := tx.Create(&w).Error; err != nil {
				return err
			}
		} else if wErr != nil {
			return wErr
		}

		outBal = newBal
		return nil
	})

	if err != nil {
		if errors.Is(err, errInsufficient) {
			return &model.ProviderCallbackResponse{Status: model.ProvCallbackInsufficientFunds, Msg: "INSUFFICIENT_BALANCE"}
		}
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INTERNAL_ERROR"}
	}
	return &model.ProviderCallbackResponse{Status: model.ProvCallbackOK, Msg: "SUCCESS", Balance: &outBal}
}

var errInsufficient = errors.New("insufficient balance")

func reqChangeBalanceOK(req *model.CallbackRequest) error {
	if strings.TrimSpace(req.Token) == "" || strings.TrimSpace(req.UserCode) == "" ||
		strings.TrimSpace(req.CurrencyCode) == "" || strings.TrimSpace(req.VendorCode) == "" {
		return errors.New("missing fields")
	}
	if req.TxnType == nil || req.WagerID == nil || strings.TrimSpace(req.TxnCode) == "" || req.Amount == nil {
		return errors.New("missing txn fields")
	}
	if req.CreatedOn == "" || req.IsFinished == nil || req.IsFreeRound == nil {
		return errors.New("missing meta fields")
	}
	return nil
}

func (s *CasinoService) callbackUpdateDetail(ctx context.Context, req *model.CallbackRequest) *model.ProviderCallbackResponse {
	if strings.TrimSpace(req.Token) == "" || req.WagerID == nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INVALID_REQUEST"}
	}
	// Token must match some api list (any currency row using this integration token)
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.APIList{}).Where("token = ?", strings.TrimSpace(req.Token)).Count(&count).Error; err != nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INTERNAL_ERROR"}
	}
	if count == 0 {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackInvalidToken, Msg: "INVALID_TOKEN"}
	}
	res := s.db.WithContext(ctx).Model(&model.CasinoWager{}).
		Where("wager_id = ?", *req.WagerID).
		Update("detail", req.Detail)
	if res.Error != nil {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackBadRequest, Msg: "INTERNAL_ERROR"}
	}
	if res.RowsAffected == 0 {
		return &model.ProviderCallbackResponse{Status: model.ProvCallbackInvalidWager, Msg: "INVALID_WAGER"}
	}
	return &model.ProviderCallbackResponse{Status: model.ProvCallbackOK, Msg: "SUCCESS"}
}

func (s *CasinoService) findAPIListByTokenAndCurrency(token, currencyCode string) (*model.APIList, error) {
	var row model.APIList
	cc := strings.ToUpper(strings.TrimSpace(currencyCode))
	err := s.db.Where("token = ? AND UPPER(TRIM(currency)) = ?", strings.TrimSpace(token), cc).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// resolvePlayerFromCallbackUserCode expects the same composite as GetGameURL: "{siteCode}_{agentCode}".
func (s *CasinoService) resolvePlayerFromCallbackUserCode(full string) (*model.User, *model.CasinoPlayer, error) {
	full = strings.TrimSpace(full)
	i := strings.LastIndex(full, "_")
	if i <= 0 || i >= len(full)-1 {
		return nil, nil, errors.New("invalid userCode")
	}
	siteCode, agentCode := full[:i], full[i+1:]
	if siteCode == "" || agentCode == "" {
		return nil, nil, errors.New("invalid userCode")
	}
	var agent model.User
	if err := s.db.Where("code = ? AND role = ?", agentCode, model.RoleAgent).First(&agent).Error; err != nil {
		return nil, nil, err
	}
	var player model.CasinoPlayer
	err := s.db.Where("agent_id = ? AND site_code = ?", agent.ID, siteCode).First(&player).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		player = model.CasinoPlayer{AgentID: agent.ID, SiteCode: siteCode}
		if err := s.db.Create(&player).Error; err != nil {
			return nil, nil, err
		}
		return &agent, &player, nil
	}
	if err != nil {
		return nil, nil, err
	}
	return &agent, &player, nil
}

func (s *CasinoService) getPlayerBalance(playerID uint, currency string) (float64, error) {
	var row model.CasinoPlayerBalance
	err := s.db.Where("player_id = ? AND currency_code = ?", playerID, currency).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.Balance, nil
}

func (s *CasinoService) lockedPlayerBalance(tx *gorm.DB, playerID uint, currency string) (*model.CasinoPlayerBalance, error) {
	var row model.CasinoPlayerBalance
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("player_id = ? AND currency_code = ?", playerID, currency).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = model.CasinoPlayerBalance{
			PlayerID:     playerID,
			CurrencyCode: currency,
			Balance:      0,
		}
		if err := tx.Create(&row).Error; err != nil {
			return nil, err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("player_id = ? AND currency_code = ?", playerID, currency).
			First(&row).Error; err != nil {
			return nil, err
		}
		return &row, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
