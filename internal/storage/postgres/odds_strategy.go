package postgres

import (
	"context"

	"github.com/pitshifer/oddshub/internal/domain"
)

type oddsStrategy interface {
	SaveOdds(ctx context.Context, provider string, odds []domain.EventOdds) error
	GetOdds(ctx context.Context, sport string) ([]domain.EventOdds, error)
}
