// Package usecase contains the business logic for price processing.
package usecase

import (
	"context"

	"github.com/ar-mokhtari/market-tracker/entity"
)

type PriceUseCase struct {
	repo    Repo
	apiKey  string
	baseURL string
}

func NewPriceUseCase(repo Repo, apiKey string, baseURL string) *PriceUseCase {
	return &PriceUseCase{
		repo:    repo,
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (uc *PriceUseCase) GetPrices(pType string) ([]entity.Price, error) {
	return uc.repo.List(pType)
}

func (uc *PriceUseCase) GetSymbolTimeline(symbol string) ([]entity.Price, error) {
	const defaultLimit = 24 // Last 24 records for hourly timeline
	return uc.repo.GetHistory(symbol, defaultLimit)
}

func (u *PriceUseCase) ListPrices(ctx context.Context, priceType string) ([]entity.Price, error) {
	return u.repo.GetAllPrices(ctx, priceType)
}
