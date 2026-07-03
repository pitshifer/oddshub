package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pitshifer/oddshub/internal/domain"
)

type Storage struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	odds   oddsStrategy
}

func New(ctx context.Context, connStr string, logger *slog.Logger, oddsWay string) (*Storage, error) {
	var pool *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= 5; attempt++ {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				break
			}
		}
		pool.Close()
		logger.Error("failed to connect to database", "attempts", attempt, "error", err)
		time.Sleep(time.Duration(attempt) * 2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// select odds strategy based on config
	var odds oddsStrategy
	switch oddsWay {
	case "lines":
		odds = &lineOddsStrategy{pool: pool, logger: logger}
	default:
		odds = &legacyOddsStrategy{pool: pool, logger: logger}
	}

	return &Storage{
		pool:   pool,
		logger: logger,
		odds:   odds,
	}, nil
}

func (s *Storage) SaveOdds(ctx context.Context, provider string, odds []domain.EventOdds) error {
	return s.odds.SaveOdds(ctx, provider, odds)
}

func (s *Storage) GetOdds(ctx context.Context, sport string) ([]domain.EventOdds, error) {
	return s.odds.GetOdds(ctx, sport)
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) SaveSports(ctx context.Context, sports []domain.Sport) error {
	for _, sport := range sports {
		_, err := s.pool.Exec(ctx,
			`INSERT INTO sports (key, title, group_name, description, active, has_outrights)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (key) DO UPDATE SET
				title = EXCLUDED.title,
				group_name = EXCLUDED.group_name,
				description = EXCLUDED.description,
				active = EXCLUDED.active,
				has_outrights = EXCLUDED.has_outrights`,
			sport.Key, sport.Title, sport.Group, sport.Description, sport.Active, sport.HasOutrights,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) GetSports(ctx context.Context) ([]domain.Sport, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, title, group_name, description, active, has_outrights FROM sports`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sports []domain.Sport
	for rows.Next() {
		var sport domain.Sport
		err := rows.Scan(&sport.Key, &sport.Title, &sport.Group, &sport.Description, &sport.Active, &sport.HasOutrights)
		if err != nil {
			return nil, err
		}
		sports = append(sports, sport)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sports, nil
}
