// Package server handles the HTTP server and routing logic.
package server

import (
	"net/http"
	"time"

	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
)

// RouterInit initializes all the routes for the application and returns an http.Server.
func RouterInit(handler *v1.Handler, port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/prices", handler.ListPrices)
	mux.HandleFunc("/api/v1/prices/fetch", handler.ManualFetch)
	mux.HandleFunc("/api/v1/prices/timeline", handler.GetTimeline)

	return &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
}
