package theoddsapi

import (
	"time"

	"github.com/pitshifer/oddshub/internal/service"
)

func mapOddsToDomain(dto []OddsResponse) []service.EventOdds {
	var result []service.EventOdds

	for _, e := range dto {
		startTime, _ := time.Parse(time.RFC3339, e.Commence)
		event := service.EventOdds{
			EventID:   e.ID,
			Sport:     e.SportKey,
			HomeTeam:  e.HomeTeam,
			AwayTeam:  e.AwayTeam,
			StartTime: startTime,
		}

		for _, bm := range e.Bookmakers {
			bookmaker := service.Bookmaker{
				Name: bm.Key,
			}

			for _, m := range bm.Markets {
				market := service.Market{
					Type: m.Key,
				}

				for _, o := range m.Outcomes {
					market.Outcomes = append(market.Outcomes, service.Outcome{
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

func mapSportsToDomain(dto []SportsResponse) []service.Sport {
	var result []service.Sport

	for _, s := range dto {
		sport := service.Sport{
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
