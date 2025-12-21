// Package v1 provides the HTTP handlers for the market tracker.
package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ar-mokhtari/market-tracker/usecase"
)

type Handler struct {
	uc *usecase.PriceUseCase
}

func NewHandler(uc *usecase.PriceUseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ListPrices(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	prices, err := h.uc.GetPrices(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			log.Printf("Encoding error: %v", encErr)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": prices}); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

func (h *Handler) ManualFetch(w http.ResponseWriter, r *http.Request) {
	err := h.uc.FetchFromExternal()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = w.Write([]byte(`{"status": "fetching successful"}`))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"timeline": timeline}); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}
