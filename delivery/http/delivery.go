// Package delivery initializes the HTTP handlers and background workers.
package delivery

import (
	"database/sql"
	"net/http"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	config "github.com/ar-mokhtari/market-tracker/config"
	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/entity"
	"github.com/ar-mokhtari/market-tracker/usecase"
)

func Init(db *sql.DB, cfg *config.Config, hub *v1.Hub) *http.ServeMux {
	repo := mysql.NewRepository(db)

	// Pass FetchInterval directly from cfg to NewPriceUseCase
	uc := usecase.NewPriceUseCase(repo, cfg.APIKey, cfg.BaseURL, cfg.FetchInterval)

	// Setup Callback to push data to WebSocket Hub
	uc.OnUpdate = func(prices []entity.Price) {
		hub.BroadcastUpdate(map[string]interface{}{"data": prices})
	}

	// Start the single worker in background
	go uc.StartAutomation()

	h := v1.NewPriceHandler(uc, hub)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return mux
}
