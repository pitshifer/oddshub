package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pitshifer/oddshub/internal/service"
)

type Handler struct {
	storage service.Storage
	client  service.Odds
}

func New(storage service.Storage, client service.Odds) *Handler {
	return &Handler{
		storage: storage,
		client:  client,
	}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/collect-sports", h.CollectSports)

		r.Route("/sports", func(r chi.Router) {
			r.Get("/{sport}/odds", h.GetOdds)
			r.Post("/{sport}/collect", h.CollectOdds)
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
	if odds == nil {
		odds = []service.EventOdds{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(odds); err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
	}
}

func (h *Handler) CollectOdds(w http.ResponseWriter, r *http.Request) {
	sport := chi.URLParam(r, "sport")
	odds, err := h.client.GetOdds(r.Context(), sport)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	err = h.storage.SaveOdds(r.Context(), "theoddsapi", odds)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}

func (h *Handler) CollectSports(w http.ResponseWriter, r *http.Request) {
	sports, err := h.client.GetSports(r.Context())
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	err = h.storage.SaveSports(r.Context(), sports)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}
