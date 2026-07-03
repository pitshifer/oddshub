package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pitshifer/oddshub/internal/domain"
)

type lineOddsStrategy struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (s *lineOddsStrategy) SaveOdds(ctx context.Context, provider string, odds []domain.EventOdds) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// upsert provider
	var providerID int
	err = tx.QueryRow(ctx,
		`INSERT INTO providers (key, name) VALUES ($1, $1)
		ON CONFLICT (key) DO UPDATE SET
			key = EXCLUDED.key
		RETURNING id`,
		provider,
	).Scan(&providerID)
	if err != nil {
		return err
	}

	sports := make(map[string]int)
	rows, err := tx.Query(ctx, "SELECT id, key FROM sports")
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		var key string
		if err := rows.Scan(&id, &key); err != nil {
			return err
		}
		sports[key] = id
	}
	rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	// upsert events
	var eventID int
	for _, event := range odds {
		sportID, ok := sports[event.Sport]
		if !ok {
			s.logger.Warn("sport not found in database", "sport", event.Sport, "event_id", event.EventID)
			continue
		}

		err := tx.QueryRow(ctx,
			`INSERT INTO events(external_id, provider_id, sport_id, home_team, away_team, start_time)
			VALUES($1, $2, $3, $4, $5, $6)
			ON CONFLICT(external_id, provider_id) DO UPDATE SET
				sport_id = EXCLUDED.sport_id,
				home_team = EXCLUDED.home_team,
				away_team = EXCLUDED.away_team,
				start_time = EXCLUDED.start_time
			RETURNING id`,
			event.EventID, providerID, sportID, event.HomeTeam, event.AwayTeam, event.StartTime,
		).Scan(&eventID)
		if err != nil {
			return err
		}

		// upsert bookmakers
		var bookmakerID int
		for _, bookmaker := range event.Bookmakers {
			err := tx.QueryRow(ctx,
				`INSERT INTO bookmakers (key)
				VALUES ($1)
				ON CONFLICT (key) DO UPDATE SET
					key = EXCLUDED.key
				RETURNING id`,
				bookmaker.Name,
			).Scan(&bookmakerID)
			if err != nil {
				return err
			}

			// upsert odds_line
			for _, market := range bookmaker.Markets {
				for _, outcome := range market.Outcomes {
					// upsert odds_line
					var lineID int
					err := tx.QueryRow(ctx,
						`INSERT INTO odds_lines (event_id, bookmaker_id, market, outcome, price, updated_at)
						VALUES ($1, $2, $3, $4, $5, NOW())
						ON CONFLICT(event_id, bookmaker_id, market, outcome) DO UPDATE SET
							price = EXCLUDED.price,
							updated_at = NOW()
						WHERE odds_lines.price IS DISTINCT FROM EXCLUDED.price
						RETURNING id`,
						eventID, bookmakerID, market.Type, outcome.Name, outcome.Price,
					).Scan(&lineID)
					if errors.Is(err, pgx.ErrNoRows) {
						continue
					}
					if err != nil {
						return err
					}

					// insert new price
					_, err = tx.Exec(ctx,
						`INSERT INTO odds_price_history (line_id, price, collected_at)
						VALUES ($1, $2, NOW())`,
						lineID, outcome.Price)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *lineOddsStrategy) GetOdds(ctx context.Context, sport string) ([]domain.EventOdds, error) {
	// Implementation for getting line odds
	return nil, nil
}
