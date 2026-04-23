package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pitshifer/oddshub/internal/service"
)

type Handler struct {
	storage service.Storage
}

func New(storage service.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Route("/sports", func(r chi.Router) {
			r.Get("/{sport}/odds", h.GetOdds)
		})
	})

	return r
}

func (h *Handler) GetOdds(w http.ResponseWriter, r *http.Request) {
	sport := chi.URLParam(r, "sport")

	odds, err := h.storage.GetOdds(r.Context(), sport)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(odds); err != nil {
		http.Error(w, "itenal service error", http.StatusInternalServerError)
	}
}
