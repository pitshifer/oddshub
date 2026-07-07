package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pitshifer/oddshub/internal/domain"
)

type OddsCollector interface {
	CollectOdds(ctx context.Context, sport string) error
	GetOdds(ctx context.Context, sport string) ([]domain.EventOdds, error)
}

type SportsCollector interface {
	CollectSports(ctx context.Context) error
	GetSports(ctx context.Context) ([]domain.Sport, error)
}

type Handler struct {
	oddsCollector   OddsCollector
	sportsCollector SportsCollector

	logger *slog.Logger
}

func New(oddsCollector OddsCollector, sportsCollector SportsCollector, logger *slog.Logger) *Handler {
	return &Handler{
		oddsCollector:   oddsCollector,
		sportsCollector: sportsCollector,
		logger:          logger,
	}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(LoggerMiddleware(h.logger))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/collect-sports", h.CollectSports)

		r.Route("/sports", func(r chi.Router) {
			r.Get("/", h.GetSports)
			r.Get("/{sport}/odds", h.GetOdds)
			r.Post("/{sport}/collect", h.CollectOdds)
		})
	})

	return r
}

func (h *Handler) GetSports(w http.ResponseWriter, r *http.Request) {
	sports, err := h.sportsCollector.GetSports(r.Context())
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(sports); err != nil {
		h.logger.Error("Failed to encode sports response", "error", err)
		http.Error(w, "internal service error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetOdds(w http.ResponseWriter, r *http.Request) {
	sport := chi.URLParam(r, "sport")

	odds, err := h.oddsCollector.GetOdds(r.Context(), sport)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}
	if odds == nil {
		odds = []domain.EventOdds{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(odds); err != nil {
		h.logger.Error("Failed to encode odds response", "error", err)
		http.Error(w, "internal service error", http.StatusInternalServerError)
	}
}

func (h *Handler) CollectOdds(w http.ResponseWriter, r *http.Request) {
	sport := chi.URLParam(r, "sport")
	err := h.oddsCollector.CollectOdds(r.Context(), sport)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, `{"status":"ok"}`)
}

func (h *Handler) CollectSports(w http.ResponseWriter, r *http.Request) {
	err := h.sportsCollector.CollectSports(r.Context())
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}
