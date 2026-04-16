package theoddsapi

import (
	"time"

	"github.com/pitshifer/oddshub/internal/service"
)

func mapToDomain(dto []OddsResponse) []service.EventOdds {
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
