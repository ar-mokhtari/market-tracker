// Package usecase contains the business logic for price processing.
package usecase

import (
	"fmt"
	"time"
)

// StartAutomation starts the background worker for periodic price fetching.
func (uc *PriceUseCase) StartAutomation() {
	// Define a helper to avoid code duplication
	runFetch := func() {
		if err := uc.fetchFromExternal(); err != nil {
			fmt.Printf("Fetch error: %v\n", err)
		}
	}

	// Run immediately on startup
	runFetch()

	ticker := time.NewTicker(time.Duration(uc.FetchInterval) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		runFetch()
	}
}
