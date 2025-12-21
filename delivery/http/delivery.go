// Package delivery initializes the HTTP handlers and background workers.
package delivery

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	config "github.com/ar-mokhtari/market-tracker/config"
	handler "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/usecase"
)

func Init(db *sql.DB, cfg *config.Config) *http.ServeMux {
	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, cfg.APIKey, cfg.BaseURL)

	// Single managed worker
	go func() {
		// Initial fetch
		_ = uc.FetchFromExternal()

		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			_ = uc.FetchFromExternal()
		}
	}()

	h := handler.NewPriceHandler(uc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return mux
}
