package gohdl

import (
	"go.gllm.dev/vanity-go/internal/services/gosvc"
	"log/slog"
	"net/http"
	"strings"
)

// Handler handles HTTP requests for Go vanity imports.
// It uses the gosvc.Service to generate the appropriate HTML responses
// for Go's import path resolution mechanism.
type Handler struct {
	service *gosvc.Service
}

// New creates a new Handler instance with the provided gosvc.Service.
// The service is responsible for generating the HTML content with proper meta tags.
func New(service *gosvc.Service) *Handler {
	return &Handler{service: service}
}

// Handle processes HTTP requests for vanity import paths.
// It expects requests in the format "/go/module/path" and generates HTML responses
// with the appropriate go-import and go-source meta tags.
//
// The handler:
//   - Returns 404 if no module path is provided
//   - Generates HTML with meta tags using the service
//   - Parses and executes the HTML template
//   - Sets proper Content-Type header
//   - Returns 404 on template execution errors
//
// Example:
//
//	Request to "/myproject" generates HTML that tells `go get` where to find
//	the actual repository for "domain.com/myproject".
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	html := h.service.Vanity(r.Context(), path)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(html))
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to write template", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
