// Package delivery initializes the HTTP handlers and background workers.
package delivery

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	config "github.com/ar-mokhtari/market-tracker/config"
	handler "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/usecase"
)

func Init(db *sql.DB, cfg *config.Config, hub *v1.Hub) *http.ServeMux {
	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, cfg.APIKey, cfg.BaseURL)
	interval := time.Duration(cfg.FetchInterval) * time.Minute
	ticker := time.NewTicker(interval)

	// Inside the background worker in delivery.Init:
	go func() {
		// Fetch and Broadcast
		if err := uc.FetchFromExternal(); err == nil {
			prices, _ := uc.GetPrices("")
			hub.BroadcastUpdate(prices)
		}

		// 2. Inside delivery.Init -> Ticker Loop
		for range ticker.C {
			if err := uc.FetchFromExternal(); err == nil {
				prices, _ := uc.GetPrices("")
				hub.BroadcastUpdate(map[string]interface{}{"data": prices})
			}
		}
	}()

	// Single managed worker
	go func() {
		// Initial fetch
		_ = uc.FetchFromExternal()

		for range ticker.C {
			_ = uc.FetchFromExternal()
		}
	}()

	h := handler.NewPriceHandler(uc, hub)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return mux
}
