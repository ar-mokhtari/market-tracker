// Package usecase implements the business logic for fetching real-time market data.
package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ar-mokhtari/market-tracker/entity"
)

// FetchFromExternal replaces the old file-reading logic with a real API call.
func (uc *PriceUseCase) FetchFromExternal() error {
	apiKey := os.Getenv("API_KEY")
	baseURL := "https://brsapi.ir/Api/Market/Gold_Currency.php?key=" + apiKey

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("api execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	// Dynamic mapping for BRS API response structure
	var result struct {
		Gold []entity.Price `json:"gold"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("json decode error: %w", err)
	}

	for _, p := range result.Gold {
		p.Type = "gold"
		// This will update the price and history ONLY if it's new
		if err := uc.repo.Upsert(p); err != nil {
			fmt.Printf("Error during upsert for %s: %v\n", p.Symbol, err)
		}
	}

	return nil
}
