package healthzhdl

import (
	"encoding/json"
	"net/http"
)

// Handler is the HTTP handler for health checks
type Handler struct{}

// New creates a new instance of the Handler for health checks
func New() *Handler {
	return &Handler{}
}

// Status represents the health check response structure
type Status struct {
	Status string `json:"status"`
}

// Healthz handles the health check endpoint
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	status := Status{
		Status: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(status)
}
