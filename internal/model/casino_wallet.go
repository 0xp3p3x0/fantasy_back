package model

import "time"

// CasinoPlayer is a site player under an agent (site user code + agent).
type CasinoPlayer struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	AgentID  uint   `json:"agent_id" gorm:"not null;uniqueIndex:ux_casino_player_agent_site"`
	SiteCode string `json:"site_code" gorm:"size:191;not null;uniqueIndex:ux_casino_player_agent_site"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CasinoPlayerBalance is per-player, per-currency wallet.
type CasinoPlayerBalance struct {
	PlayerID     uint    `json:"player_id" gorm:"primaryKey"`
	CurrencyCode string  `json:"currency_code" gorm:"primaryKey;size:16"`
	Balance      float64 `json:"balance" gorm:"default:0"`
	UpdatedAt    time.Time
}

// CasinoWager stores provider wager metadata for UpdateDetail and lookups.
type CasinoWager struct {
	WagerID      int64  `json:"wager_id" gorm:"primaryKey"`
	PlayerID     uint   `json:"player_id" gorm:"not null;index"`
	AgentID      uint   `json:"agent_id" gorm:"not null;index"`
	CurrencyCode string `json:"currency_code" gorm:"size:16;not null"`
	VendorCode   string `json:"vendor_code" gorm:"size:128"`
	GameCode     string `json:"game_code" gorm:"size:128"`
	GameRoundID  string `json:"game_round_id" gorm:"size:128"`
	Detail       string `json:"detail" gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CasinoProcessedTxn ensures ChangeBalance is idempotent per txnCode.
type CasinoProcessedTxn struct {
	TxnCode      string    `json:"txn_code" gorm:"primaryKey;size:128"`
	WagerID      int64     `json:"wager_id" gorm:"not null;index"`
	PlayerID     uint      `json:"player_id" gorm:"not null"`
	CurrencyCode string    `json:"currency_code" gorm:"size:16;not null"`
	BalanceAfter float64   `json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}
