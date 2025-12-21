// Package usecase contains the business logic for price processing.
package usecase

import (
	"fmt"
	"time"
)

// StartAutomation starts the background worker for periodic price fetching.
func (uc *PriceUseCase) StartAutomation() {
	// Immediate first fetch
	if err := uc.FetchFromExternal(); err != nil {
		fmt.Printf("Initial automation fetch failed: %v\n", err)
	}

	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("Executing scheduled fetch at: %v\n", time.Now().Format("15:04:05"))
		if err := uc.FetchFromExternal(); err != nil {
			fmt.Printf("Automated fetch error: %v\n", err)
		}
	}
}
