package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pitshifer/oddshub/internal/service"
)

type legacyOddsStrategy struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (s *legacyOddsStrategy) SaveOdds(ctx context.Context, provider string, odds []service.EventOdds) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Upsert provider
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
	rows, err := tx.Query(ctx, `SELECT id, key FROM sports`)
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

	// Upsert events
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

		// Upsert bookmakers
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

			// Upsert odds
			for _, market := range bookmaker.Markets {
				for _, outcome := range market.Outcomes {
					_, err := tx.Exec(ctx,
						`INSERT INTO odds (event_id, bookmaker_id, market, outcome, price)
						VALUES ($1, $2, $3, $4, $5)`,
						eventID, bookmakerID, market.Type, outcome.Name, outcome.Price,
					)
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

func (s *legacyOddsStrategy) GetOdds(ctx context.Context, sport string) ([]service.EventOdds, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT e.external_id, s.key, e.home_team, e.away_team, e.start_time, b.key, o.market, o.outcome, o.price
		FROM events e
		JOIN odds o ON e.id = o.event_id
		JOIN bookmakers b ON o.bookmaker_id = b.id
		JOIN sports s ON e.sport_id = s.id
		WHERE s.key = $1
		ORDER BY e.id, b.key, o.market`,
		sport,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []service.EventOdds
	eventsMap := make(map[string]int)
	bookmakerMap := make(map[string]int)
	marketMap := make(map[string]int)

	for rows.Next() {
		var (
			eventID                                               string
			sport, homeTeam, awayTeam, bookmaker, market, outcome string
			startTime                                             time.Time
			price                                                 float64
		)

		err := rows.Scan(&eventID, &sport, &homeTeam, &awayTeam, &startTime, &bookmaker, &market, &outcome, &price)
		if err != nil {
			return nil, err
		}

		ei, ok := eventsMap[eventID]
		if !ok {
			eventOdds := service.EventOdds{
				EventID:   eventID,
				Sport:     sport,
				HomeTeam:  homeTeam,
				AwayTeam:  awayTeam,
				StartTime: startTime,
			}
			result = append(result, eventOdds)
			ei = len(result) - 1
			eventsMap[eventID] = ei
		}

		bmKey := fmt.Sprintf("%s-%s", eventID, bookmaker)
		bi, ok := bookmakerMap[bmKey]
		if !ok {
			bookmakerOdds := service.Bookmaker{
				Name: bookmaker,
			}
			result[ei].Bookmakers = append(result[ei].Bookmakers, bookmakerOdds)
			bi = len(result[ei].Bookmakers) - 1
			bookmakerMap[bmKey] = bi
		}

		mKey := fmt.Sprintf("%s-%s", bmKey, market)
		mi, ok := marketMap[mKey]
		if !ok {
			marketOdds := service.Market{
				Type: market,
			}
			result[ei].Bookmakers[bi].Markets = append(result[ei].Bookmakers[bi].Markets, marketOdds)
			mi = len(result[ei].Bookmakers[bi].Markets) - 1
			marketMap[mKey] = mi
		}

		outcomeOdds := service.Outcome{
			Name:  outcome,
			Price: price,
		}
		result[ei].Bookmakers[bi].Markets[mi].Outcomes = append(result[ei].Bookmakers[bi].Markets[mi].Outcomes, outcomeOdds)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
