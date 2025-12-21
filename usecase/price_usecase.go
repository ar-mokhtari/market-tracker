// Package usecase contains the business logic for price processing.
package usecase

import (
	"os"

	"github.com/ar-mokhtari/market-tracker/entity"
)

type PriceUseCase struct {
	repo    Repo
	apiKey  string
	baseURL string
}

func NewPriceUseCase(repo Repo, apiKey string) *PriceUseCase {
	baseUrl := os.Getenv("API_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://BrsApi.ir/Api/Market/Gold_Currency.php"
	}
	return &PriceUseCase{repo: repo, apiKey: apiKey, baseURL: baseUrl}
}

func (uc *PriceUseCase) GetPrices(pType string) ([]entity.Price, error) {
	return uc.repo.List(pType)
}

func (uc *PriceUseCase) GetSymbolTimeline(symbol string) ([]entity.Price, error) {
	const defaultLimit = 24 // Last 24 records for hourly timeline
	return uc.repo.GetHistory(symbol, defaultLimit)
}
