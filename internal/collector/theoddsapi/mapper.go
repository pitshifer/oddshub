package theoddsapi

import (
	"time"

	"github.com/pitshifer/oddshub/internal/domain"
)

func mapOddsToDomain(dto []OddsResponse) []domain.EventOdds {
	var result []domain.EventOdds

	for _, e := range dto {
		startTime, _ := time.Parse(time.RFC3339, e.Commence)
		event := domain.EventOdds{
			EventID:   e.ID,
			Sport:     e.SportKey,
			HomeTeam:  e.HomeTeam,
			AwayTeam:  e.AwayTeam,
			StartTime: startTime,
		}

		for _, bm := range e.Bookmakers {
			bookmaker := domain.Bookmaker{
				Name: bm.Key,
			}

			for _, m := range bm.Markets {
				market := domain.Market{
					Type: m.Key,
				}

				for _, o := range m.Outcomes {
					market.Outcomes = append(market.Outcomes, domain.Outcome{
						Name:  o.Name,
						Price: o.Price,
					})
				}

				bookmaker.Markets = append(bookmaker.Markets, market)
			}

			event.Bookmakers = append(event.Bookmakers, bookmaker)
		}

		result = append(result, event)
	}

	return result
}

func mapSportsToDomain(dto []SportsResponse) []domain.Sport {
	var result []domain.Sport

	for _, s := range dto {
		sport := domain.Sport{
			Key:          s.Key,
			Title:        s.Title,
			Group:        s.Group,
			Description:  s.Description,
			Active:       s.Active,
			HasOutrights: s.HasOutrights,
		}

		result = append(result, sport)
	}

	return result
}
