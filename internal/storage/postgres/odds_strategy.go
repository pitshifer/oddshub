package postgres

import (
	"context"

	"github.com/pitshifer/oddshub/internal/service"
)

type oddsStrategy interface {
	SaveOdds(ctx context.Context, provider string, odds []service.EventOdds) error
	GetOdds(ctx context.Context, sport string) ([]service.EventOdds, error)
}
