package application

import (
	"context"
	"log/slog"

	"github.com/pitshifer/oddshub/internal/domain"
)

type SportsService struct {
	client  domain.Odds
	storage domain.Storage
	logger  *slog.Logger
}

func NewSportsService(storage domain.Storage, client domain.Odds, logger *slog.Logger) *SportsService {
	return &SportsService{
		storage: storage,
		client:  client,
		logger:  logger,
	}
}

func (s *SportsService) CollectSports(ctx context.Context) error {
	sports, err := s.client.GetSports(ctx)
	if err != nil {
		s.logger.Error("Failed to collect sports", "error", err)
		return err
	}

	err = s.storage.SaveSports(ctx, sports)
	if err != nil {
		s.logger.Error("Failed to save sports", "error", err)
		return err
	}

	return nil
}

func (s *SportsService) GetSports(ctx context.Context) ([]domain.Sport, error) {
	sports, err := s.storage.GetSports(ctx)
	if err != nil {
		s.logger.Error("Failed to get sports", "error", err)
		return nil, err
	}

	return sports, nil
}
