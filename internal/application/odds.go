package application

import (
	"context"
	"log/slog"

	"github.com/pitshifer/oddshub/internal/domain"
)

type OddsService struct {
	client  domain.Odds
	storage domain.Storage
	logger  *slog.Logger
}

func NewOddsService(storage domain.Storage, client domain.Odds, logger *slog.Logger) *OddsService {
	return &OddsService{
		storage: storage,
		client:  client,
		logger:  logger,
	}
}

func (o *OddsService) CollectOdds(ctx context.Context, sport string) error {
	odds, err := o.client.GetOdds(ctx, sport)
	if err != nil {
		o.logger.Error("Failed to collect odds for sport", "sport", sport, "error", err)
		return err
	}

	err = o.storage.SaveOdds(ctx, "theoddsapi", odds)
	if err != nil {
		o.logger.Error("Failed to save odds for sport", "sport", sport, "error", err)
		return err
	}

	return nil
}

func (o *OddsService) GetOdds(ctx context.Context, sport string) ([]domain.EventOdds, error) {
	odds, err := o.storage.GetOdds(ctx, sport)
	if err != nil {
		o.logger.Error("Failed to get odds for sport", "sport", sport, "error", err)
		return nil, err
	}

	return odds, nil
}
