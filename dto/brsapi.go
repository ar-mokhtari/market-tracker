package dto

// PriceResponse defines how price data looks like for the client
type PriceResponse struct {
	Date   string  `json:"date"`
	Time   string  `json:"time"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Type   string  `json:"type"`
	Unit   string  `json:"unit"`
}
