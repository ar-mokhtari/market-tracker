// Package delivery initializes the HTTP handlers and background workers.
package delivery

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	config "github.com/ar-mokhtari/market-tracker/config"
	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/entity"
	"github.com/ar-mokhtari/market-tracker/usecase"
)

func Init(db *sql.DB, cfg *config.Config, hub *v1.Hub) *http.ServeMux {
	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, cfg.APIKey, cfg.BaseURL)

	// Setup Callback
	uc.OnUpdate = func(prices []entity.Price) {
		hub.BroadcastUpdate(map[string]interface{}{"data": prices})
	}

	// Single Optimized Worker
	go func() {
		// Ensure interval is at least 1 minute if config fails
		intervalMin := cfg.FetchInterval
		if intervalMin <= 0 {
			intervalMin = 1
		}

		interval := time.Duration(intervalMin) * time.Minute
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Printf("Worker started. Interval: %v", interval)

		// First execution
		if err := uc.FetchFromExternal(); err != nil {
			log.Printf("Initial fetch error: %v", err)
		}

		for {
			select {
			case t := <-ticker.C:
				log.Printf("Ticker fired at %v. Fetching data...", t.Format("15:04:05"))
				if err := uc.FetchFromExternal(); err != nil {
					log.Printf("Scheduled fetch error: %v", err)
				}
			}
		}
	}()

	h := v1.NewPriceHandler(uc, hub)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}
