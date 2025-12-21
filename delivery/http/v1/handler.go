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
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": prices})
}

func (h *Handler) ManualFetch(w http.ResponseWriter, r *http.Request) {
	// حذف کلمه go برای دیدن خطاها
	err := h.uc.FetchFromExternal()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "fetching successful"}`))
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
