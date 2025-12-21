// Package v1 main provides the entry point for the market-tracker service,
// handling dependency injection, database connection, and background workers.
package v1

import (
	"encoding/json"
	"net/http"

	"github.com/ar-mokhtari/market-tracker/usecase"
)

type Handler struct {
	uc *usecase.PriceUseCase
}

func NewPriceHandler(uc *usecase.PriceUseCase) *Handler {
	return &Handler{uc: uc}
}

// GetPrices maps to /api/v1/prices
func (h *Handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	prices, err := h.uc.GetPrices(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": prices})
}

// FetchPrices maps to /api/v1/prices/fetch
func (h *Handler) FetchPrices(w http.ResponseWriter, r *http.Request) {
	err := h.uc.FetchFromExternal()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, _ = w.Write([]byte(`{"status": "fetching successful"}`))
}

// GetTimeline maps to /api/v1/prices/timeline
func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}

	timeline, err := h.uc.GetSymbolTimeline(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"timeline": timeline})
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/prices", h.GetPrices)
	mux.HandleFunc("/api/v1/prices/fetch", h.FetchPrices)
	mux.HandleFunc("/api/v1/prices/timeline", h.GetTimeline)
}
