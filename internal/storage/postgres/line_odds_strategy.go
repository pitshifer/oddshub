package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pitshifer/oddshub/internal/service"
)

type lineOddsStrategy struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (s *lineOddsStrategy) SaveOdds(ctx context.Context, provider string, odds []service.EventOdds) error {
	// Implementation for saving line odds
	return nil
}

func (s *lineOddsStrategy) GetOdds(ctx context.Context, sport string) ([]service.EventOdds, error) {
	// Implementation for getting line odds
	return nil, nil
}
