// Package usecase contains the business logic for price processing.
package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/ar-mokhtari/market-tracker/entity"
)

type PriceUseCase struct {
	repo          Repo
	apiKey        string
	baseURL       string
	httpClient    *http.Client
	OnUpdate      func([]entity.Price)
	FetchInterval int
}

func NewPriceUseCase(repo Repo, apiKey string, baseURL string, interval int) *PriceUseCase {
	return &PriceUseCase{
		repo:          repo,
		apiKey:        apiKey,
		baseURL:       baseURL,
		FetchInterval: interval,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (uc *PriceUseCase) GetPrices(pType string) ([]entity.Price, error) {
	return uc.repo.List(pType)
}

func (uc *PriceUseCase) GetSymbolTimeline(symbol string) ([]entity.Price, error) {
	const defaultLimit = 24 // Last 24 records for hourly timeline
	return uc.repo.GetHistory(symbol, defaultLimit)
}

func (uc *PriceUseCase) ListPrices(ctx context.Context, priceType string) ([]entity.Price, error) {
	return uc.repo.GetAllPrices(ctx, priceType)
}
