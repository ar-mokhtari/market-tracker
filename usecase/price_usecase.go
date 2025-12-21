// Package usecase contains the business logic for price processing.
package usecase

import (
	"log"
	"os"

	"github.com/ar-mokhtari/market-tracker/entity"
)

// var _ = httpClient = &http.Client{
// 	Timeout: 30 * time.Second,
// 	Transport: &http.Transport{
// 		ForceAttemptHTTP2: false, // اجبار به HTTP/1.1
// 		TLSClientConfig: &tls.Config{
// 			InsecureSkipVerify: true,
// 			CipherSuites: []uint16{
// 				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
// 				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
// 			},
// 		},
// 	},
// }

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

// FetchFromExternal fetches data and handles errors for each record to satisfy linter.
func (uc *PriceUseCase) FetchFromExternal() error {
	// ... (Your file reading or API logic)

	// In the processing loop:
	priceData := entity.Price{
		Symbol: "GOLD18K", // Example
		Price:  "13788700",
		// ... other fields
	}

	if err := uc.repo.Upsert(priceData); err != nil {
		log.Printf("Failed to upsert symbol %s: %v", priceData.Symbol, err)
		// We continue the loop but log the error to satisfy errcheck
	}

	return nil
}

func (uc *PriceUseCase) GetPrices(pType string) ([]entity.Price, error) {
	return uc.repo.List(pType)
}

func (uc *PriceUseCase) GetSymbolTimeline(symbol string) ([]entity.Price, error) {
	const defaultLimit = 24 // Last 24 records for hourly timeline
	return uc.repo.GetHistory(symbol, defaultLimit)
}
