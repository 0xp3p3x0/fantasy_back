package model

// CallbackRequest is the union body for a single provider callback URL.
// The vendor POSTs different JSON shapes to the same endpoint; "method" selects behavior:
//   GetBalance   → token, userCode, currencyCode
//   ChangeBalance → those plus vendorCode, txnType, txnCode, wagerId, amount, createdOn, isFinished, isFreeRound, etc.
//   UpdateDetail → token, wagerId, detail
// Unused fields are omitted in JSON and read as zero values here.
type CallbackRequest struct {
	Method       string   `json:"method"`
	Token        string   `json:"token"`
	UserCode     string   `json:"userCode"`
	CurrencyCode string   `json:"currencyCode"`
	VendorCode   string   `json:"vendorCode"`
	TxnType      *int     `json:"txnType"`
	TxnCode      string   `json:"txnCode"`
	PairCode     *string  `json:"pairCode"`
	WagerID      *int64   `json:"wagerId"`
	Detail       string   `json:"detail"`
	Amount       *float64 `json:"amount"`
	GameCode     string   `json:"gameCode"`
	GameRoundID  string   `json:"gameRoundId"`
	CreatedOn    string   `json:"createdOn"`
	IsFinished   *bool    `json:"isFinished"`
	IsFreeRound  *bool    `json:"isFreeRound"`
}

// ProviderCallbackResponse is the JSON shape the upstream provider expects (not APIResponse).
type ProviderCallbackResponse struct {
	Status  int      `json:"status"`
	Msg     string   `json:"msg"`
	Balance *float64 `json:"balance,omitempty"`
}

const (
	ProvCallbackOK                 = 0
	ProvCallbackInvalidToken       = 1
	ProvCallbackInvalidUser        = 5
	ProvCallbackInsufficientFunds  = 6
	ProvCallbackInvalidWager       = 18
	ProvCallbackBadRequest         = 99
)
