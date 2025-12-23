package v1

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/prices", h.GetPrices)
	mux.HandleFunc("/api/v1/prices/fetch", h.FetchPrices)
	mux.HandleFunc("/api/v1/prices/timeline", h.GetTimeline)
	mux.HandleFunc("/api/v1/prices/all", h.ListAllPrices)

	// Add WebSocket endpoint
	mux.HandleFunc("/ws", h.hub.ServeWS)

	// Docker health check endpoint
	mux.HandleFunc("/health", h.HealthCheck) // Docker health check endpoint
}
