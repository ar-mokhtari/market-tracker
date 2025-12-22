// Package v1 provides the delivery layer for price tracking services.
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
// It returns a list of the most recent prices, optionally filtered by the 'type' query parameter.
func (h *Handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	prices, err := h.uc.GetPrices(t)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": prices}); err != nil {
		return
	}
}

// FetchPrices maps to /api/v1/prices/fetch
// It triggers an external API call to fetch current market prices and store them in the database.
func (h *Handler) FetchPrices(w http.ResponseWriter, r *http.Request) {
	err := h.uc.FetchFromExternal()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	if _, err := w.Write([]byte(`{"status": "fetching successful"}`)); err != nil {
		return
	}
}

// GetTimeline maps to /api/v1/prices/timeline
// It provides a historical timeline of prices for a specific symbol required in query parameters.
func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}

	timeline, err := h.uc.GetSymbolTimeline(symbol)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"timeline": timeline}); err != nil {
		return
	}
}

// ListAllPrices maps to /api/v1/prices/all
// It retrieves all price records from the database with an optional filter by type.
func (h *Handler) ListAllPrices(w http.ResponseWriter, r *http.Request) {
	priceType := r.URL.Query().Get("type")

	prices, err := h.uc.ListPrices(r.Context(), priceType)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": prices}); err != nil {
		return
	}
}
