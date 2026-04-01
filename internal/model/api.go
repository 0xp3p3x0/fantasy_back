package model

type GetGameURLRequest struct {
	UserCode     string `json:"userCode" binding:"required"`
	Nickname     string `json:"nickname,omitempty"`
	VendorCode   string `json:"vendorCode" binding:"required"`
	GameCode     string `json:"gameCode,omitempty"`
	CurrencyCode string `json:"currencyCode" binding:"required"`
	Language     string `json:"language,omitempty"`
	Channel      string `json:"channel,omitempty"`
	IsDemo       bool   `json:"isDemo"`
}

type GetGameURLResponse struct {
	LaunchURL string `json:"launchUrl"`
}

type BalanceResponse struct {
	UserCode string             `json:"userCode"`
	Balances map[string]float64 `json:"balances"`
}

type GamesResponse struct {
	Vendors []Vendor `json:"vendors,omitempty"`
	Games   []Game   `json:"games,omitempty"`
}

type PlaceBetRequest struct {
	UserCode     string  `json:"userCode" binding:"required"`
	CurrencyCode string  `json:"currencyCode" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	Action       string  `json:"action" binding:"required"` // deposit | withdraw | withdraw_all
}

type BetResponse struct {
	PrevBalance float64 `json:"prevBalance"`
	Balance     float64 `json:"balance"`
}

type BetDetailResponse struct {
	Wager Wager `json:"wager"`
}

type Vendor struct {
	VendorCode string `json:"vendorCode"`
	VendorName string `json:"vendorName"`
	GameType   int    `json:"gameType"`
}

type Game struct {
	GameCode string `json:"gameCode"`
	GameName string `json:"gameName"`
	GameType int    `json:"gameType"`
	ImageURL string `json:"imageUrl"`
}

type Wager struct {
	UserCode      string  `json:"userCode"`
	VendorCode    string  `json:"vendorCode"`
	GameType      int     `json:"gameType"`
	GameCode      string  `json:"gameCode"`
	GameRoundID   string  `json:"gameRoundId"`
	WagerID       int64   `json:"wagerId"`
	CurrencyCode  string  `json:"currencyCode"`
	BetAmount     float64 `json:"betAmount"`
	PayoutAmount  float64 `json:"payoutAmount"`
	BeforeBalance float64 `json:"beforeBalance"`
	AfterBalance  float64 `json:"afterBalance"`
	Detail        string  `json:"detail"`
	CreatedOn     string  `json:"createdOn"`
	ModifiedOn    string  `json:"modifiedOn"`
	SettlementOn  string  `json:"settlementOn"`
	IsFinished    bool    `json:"isFinished"`
	Status        int     `json:"status"`
}
