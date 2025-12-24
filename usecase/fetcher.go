// Package usecase implements the business logic for fetching real-time market data.
package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ar-mokhtari/market-tracker/entity"
)

func (uc *PriceUseCase) FetchFromExternal() error {
	fullURL := fmt.Sprintf("%s?key=%s", uc.baseURL, uc.apiKey)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 OPR/106.0.0.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network level error: %w", err)
	}
	defer resp.Body.Close()

	var result map[string][]entity.Price
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	var allUpdated []entity.Price
	for category, prices := range result {
		for _, p := range prices {
			p.Type = category
			_ = uc.repo.Upsert(p)
			allUpdated = append(allUpdated, p)
		}
	}

	if uc.OnUpdate != nil && len(allUpdated) > 0 {
		uc.OnUpdate(allUpdated)
	}
	return nil
}
