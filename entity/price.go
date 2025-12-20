// Package entity is about fields and models of domain
package entity

import "time"

type Price struct {
	ID            uint      `json:"id"`
	Date          string    `json:"date"`
	Time          string    `json:"time"`
	TimeUnix      int64     `json:"time_unix"`
	Symbol        string    `json:"symbol"`
	NameEn        string    `json:"name_en"`
	NameFa        string    `json:"name_fa"`
	Price         string    `json:"price"`
	ChangeValue   string    `json:"change_value"`
	ChangePercent float64   `json:"change_percent"`
	Unit          string    `json:"unit"`
	Type          string    `json:"type"` // gold, currency, cryptocurrency
	MarketCap     *int64    `json:"market_cap,omitempty"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
