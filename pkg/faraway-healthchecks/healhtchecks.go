package farawayhealthchecks

import (
	"net/http"
)

// Handler returns an http.Handler for health check endpoints.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/live", liveHandler)
	mux.HandleFunc("/ready", readyHandler)
	return mux
}

// liveHandler handles the liveness probe endpoint.
func liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// readyHandler handles the readiness probe endpoint.
func readyHandler(w http.ResponseWriter, r *http.Request) {
	// Implement readiness logic here (e.g., checking external dependencies)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
