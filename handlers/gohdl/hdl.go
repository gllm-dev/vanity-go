package gohdl

import (
	"go.gllm.dev/vanity/services/gosvc"
	"html/template"
	"net/http"
	"strings"
)

type Handler struct {
	service *gosvc.Service
}

func New(service *gosvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/go/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	html := h.service.Vanity(r.Context(), path)
	tmpl := template.Must(template.New("go-import").Parse(html))
	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}
