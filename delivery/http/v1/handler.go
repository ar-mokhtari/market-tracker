// Package v1 provides the delivery layer for price tracking services.
package v1

import (
	"encoding/json"
	"net/http"

	"github.com/ar-mokhtari/market-tracker/dto"
	"github.com/ar-mokhtari/market-tracker/usecase"
)

type Handler struct {
	uc  *usecase.PriceUseCase
	hub *Hub // Add this
}

func NewPriceHandler(uc *usecase.PriceUseCase, h *Hub) *Handler {
	return &Handler{
		uc:  uc,
		hub: h,
	}
}

// sendError is a helper to return JSON error responses consistently.
func (h *Handler) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		http.Error(w, message, code)
	}
}

// GetPrices maps to /api/v1/prices
// It returns a list of the most recent prices converted to DTOs.
func (h *Handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	prices, err := h.uc.GetPrices(t)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]dto.PriceResponse, 0, len(prices))
	for _, p := range prices {
		actualPrice, _ := p.Price.Float64()
		response = append(response, dto.PriceResponse{
			Date:   p.Date,
			Time:   p.Time,
			Symbol: p.Symbol,
			Price:  actualPrice,
			Type:   p.Type,
			Unit:   p.Unit,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": response}); err != nil {
		return
	}
}

// FetchPrices maps to /api/v1/prices/fetch
// Triggering manual fetch is now deprecated to protect system resources.
// It will now just return the latest data from the database.
func (h *Handler) FetchPrices(w http.ResponseWriter, r *http.Request) {
	// Instead of uc.FetchFromExternal(), we just get latest prices from DB
	// This prevents API rate limiting and high CPU usage.
	prices, err := h.uc.GetPrices("")
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Data retrieved from local storage",
		"data":    prices,
	})
}

// GetTimeline maps to /api/v1/prices/timeline
// It provides a historical timeline of prices mapped to DTOs for a specific symbol.
func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		h.sendError(w, "symbol is required", http.StatusBadRequest)
		return
	}

	timeline, err := h.uc.GetSymbolTimeline(symbol)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]dto.PriceResponse, 0, len(timeline))
	for _, p := range timeline {
		actualPrice, _ := p.Price.Float64()
		response = append(response, dto.PriceResponse{
			Date:   p.Date,
			Time:   p.Time,
			Symbol: p.Symbol,
			Price:  actualPrice,
			Type:   p.Type,
			Unit:   p.Unit,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"timeline": response}); err != nil {
		return
	}
}

// ListAllPrices maps to /api/v1/prices/all
// It retrieves all price records and maps them to DTOs for the response.
func (h *Handler) ListAllPrices(w http.ResponseWriter, r *http.Request) {
	priceType := r.URL.Query().Get("type")

	prices, err := h.uc.ListPrices(r.Context(), priceType)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]dto.PriceResponse, 0, len(prices))
	for _, p := range prices {
		actualPrice, _ := p.Price.Float64()
		response = append(response, dto.PriceResponse{
			Date:   p.Date,
			Time:   p.Time,
			Symbol: p.Symbol,
			Price:  actualPrice,
			Type:   p.Type,
			Unit:   p.Unit,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": response}); err != nil {
		return
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respond(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *Handler) respond(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// If encoding fails, fallback to basic error
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
