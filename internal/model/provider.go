package model

type ProviderBaseResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type ProviderGetGameURLRequest struct {
	Method         string  `json:"method"`
	Token          string  `json:"token"`
	AgentCode      string  `json:"agentCode"`
	UserCode       string  `json:"userCode"`
	Nickname       string  `json:"nickname,omitempty"`
	VendorCode     string  `json:"vendorCode"`
	GameCode       string  `json:"gameCode,omitempty"`
	CurrencyCode   string  `json:"currencyCode"`
	Language       string  `json:"language,omitempty"`
	Channel        string  `json:"channel,omitempty"`
	FreeRounds     string  `json:"freeRounds,omitempty"`
	FreeRoundsCode string  `json:"freeRoundsCode,omitempty"`
	CustomGameName string  `json:"customGameName,omitempty"`
	HomeUrl        string  `json:"homeUrl,omitempty"`
	IsDemo         bool    `json:"isDemo,omitempty"`
	UserBalance    float64 `json:"userBalance,omitempty"`
}

type ProviderGetGameURLResponse struct {
	ProviderBaseResponse
	LaunchURL string `json:"launchUrl"`
}

type ProviderGetUserInfoRequest struct {
	Method    string `json:"method"`
	Token     string `json:"token"`
	AgentCode string `json:"agentCode"`
	UserCode  string `json:"userCode,omitempty"`
}

type ProviderUser struct {
	UserCode string             `json:"userCode"`
	Balances map[string]float64 `json:"balances"`
}

type ProviderGetUserInfoResponse struct {
	ProviderBaseResponse
	Users []ProviderUser `json:"users"`
}

type ProviderGetVendorsRequest struct {
	Method    string `json:"method"`
	Token     string `json:"token"`
	AgentCode string `json:"agentCode"`
}

type ProviderVendor struct {
	VendorCode string `json:"vendorCode"`
	VendorName string `json:"vendorName"`
	GameType   int    `json:"gameType"`
}

type ProviderGetVendorsResponse struct {
	ProviderBaseResponse
	Vendors []ProviderVendor `json:"vendors"`
}

type ProviderGetVendorGamesRequest struct {
	Method     string `json:"method"`
	Token      string `json:"token"`
	AgentCode  string `json:"agentCode"`
	VendorCode string `json:"vendorCode"`
}

type ProviderVendorGame struct {
	GameCode string `json:"gameCode"`
	GameName string `json:"gameName"`
	GameType int    `json:"gameType"`
	ImageURL string `json:"imageUrl"`
}

type ProviderGetVendorGamesResponse struct {
	ProviderBaseResponse
	VendorGames []ProviderVendorGame `json:"vendorGames"`
}

type ProviderChangeBalanceRequest struct {
	Method       string  `json:"method"`
	Token        string  `json:"token"`
	AgentCode    string  `json:"agentCode"`
	UserCode     string  `json:"userCode"`
	CurrencyCode string  `json:"currencyCode"`
	Amount       float64 `json:"amount,omitempty"`
}

type ProviderChangeBalanceResponse struct {
	ProviderBaseResponse
	PrevBalance float64 `json:"prevBalance"`
	Balance     float64 `json:"balance"`
}

type ProviderGetWagerInfoRequest struct {
	Method    string `json:"method"`
	Token     string `json:"token"`
	AgentCode string `json:"agentCode"`
	WagerID   int64  `json:"wagerId"`
}

type ProviderWagerResponse struct {
	ProviderBaseResponse
	Wager Wager `json:"wager"`
}
