package worker

import (
	"context"
	"log/slog"
	"time"
)

type OddsCollector interface {
	CollectOdds(ctx context.Context, sport string) error
}

type Worker struct {
	oddsCollector OddsCollector
	interval      time.Duration
	sports        []string
	logger        *slog.Logger
}

func NewWorker(oddsCollector OddsCollector, interval time.Duration, sports []string, logger *slog.Logger) *Worker {
	return &Worker{
		oddsCollector: oddsCollector,
		interval:      interval,
		sports:        sports,
		logger:        logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	t := time.NewTicker(w.interval)
	for {
		select {
		case <-t.C:
			for _, sport := range w.sports {
				w.logger.Info("Collecting odds", slog.String("sport", sport))
				err := w.oddsCollector.CollectOdds(ctx, sport)
				if err != nil {
					w.logger.Error("Failed to collect odds", slog.String("sport", sport), "err", err)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
