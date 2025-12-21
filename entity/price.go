// Package entity defines the domain models for the market tracker.
package entity

import (
	"encoding/json"
	"time"
)

// Price represents the full market data for any symbol.
type Price struct {
	ID            uint        `json:"id"`
	Date          string      `json:"date"`
	Time          string      `json:"time"`
	TimeUnix      int64       `json:"time_unix"`
	Symbol        string      `json:"symbol"`
	NameEn        string      `json:"name_en"`
	NameFa        string      `json:"name_fa"`
	Price         json.Number `json:"price"`        // Smart handling
	ChangeValue   json.Number `json:"change_value"` // Fixed: Handles numbers from API
	ChangePercent float64     `json:"change_percent"`
	Unit          string      `json:"unit"`
	Type          string      `json:"type"`
	MarketCap     *int64      `json:"market_cap,omitempty"`
	Description   string      `json:"description,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
