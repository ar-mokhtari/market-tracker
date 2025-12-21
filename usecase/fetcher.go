// Package usecase implements the business logic for fetching real-time market data.
package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ar-mokhtari/market-tracker/entity"
)

// FetchFromExternal triggers the real API call with browser-like headers.
func (uc *PriceUseCase) FetchFromExternal() error {
	fullURL := fmt.Sprintf("%s?key=%s", uc.baseURL, uc.apiKey)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set browser headers to avoid 'connection reset'
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: 20 * time.Second, // Increased timeout for slower network responses
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("api execution failed (network level): %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var result struct {
		Gold []entity.Price `json:"gold"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode api response: %w", err)
	}

	for _, p := range result.Gold {
		p.Type = "gold"
		if err := uc.repo.Upsert(p); err != nil {
			fmt.Printf("Error during upsert for %s: %v\n", p.Symbol, err)
		}
	}

	return nil
}
