package v1

import (
	"encoding/json"
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
	go h.uc.FetchFromExternal()
	w.Write([]byte(`{"status": "fetching started"}`))
}
