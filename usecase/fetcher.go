// Package usecase implements the business logic for fetching real-time market data.
package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ar-mokhtari/market-tracker/entity"
)

// FetchFromExternal triggers the real API call with browser-like headers.
func (uc *PriceUseCase) FetchFromExternal() error {
	fullURL := fmt.Sprintf("%s?key=%s", uc.baseURL, uc.apiKey)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Browser-like headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	// Optimized: Using the shared httpClient
	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network level error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var result struct {
		Gold []entity.Price `json:"gold"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	for _, p := range result.Gold {
		p.Type = "gold"
		_ = uc.repo.Upsert(p)
	}

	return nil
}
