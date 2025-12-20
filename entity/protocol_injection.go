package entity

import "context"

type PriceRepository interface {
	Create(ctx context.Context, price *Price) error
	Update(ctx context.Context, price *Price) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*Price, error)
	GetAll(ctx context.Context, priceType string) ([]Price, error)
	GetBySymbol(ctx context.Context, symbol string) (*Price, error)
}

type PriceUseCase interface {
	Create(ctx context.Context, req CreatePriceRequest) (*Price, error)
	Update(ctx context.Context, id uint, req UpdatePriceRequest) (*Price, error)
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*Price, error)
	GetAll(ctx context.Context, priceType string) ([]Price, error)
	FetchAndSaveMarketData(ctx context.Context) error
}

type CreatePriceRequest struct {
	Symbol        string  `json:"symbol"`
	NameEn        string  `json:"name_en"`
	NameFa        string  `json:"name_fa"`
	Price         string  `json:"price"`
	ChangeValue   string  `json:"change_value"`
	ChangePercent float64 `json:"change_percent"`
	Unit          string  `json:"unit"`
	Type          string  `json:"type"`
}

type UpdatePriceRequest struct {
	Price         string  `json:"price"`
	ChangeValue   string  `json:"change_value"`
	ChangePercent float64 `json:"change_percent"`
}
