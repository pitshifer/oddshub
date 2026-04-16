package theoddsapi

type OddsResponse struct {
	ID         string `json:"id"`
	SportKey   string `json:"sport_key"`
	HomeTeam   string `json:"home_team"`
	AwayTeam   string `json:"away_team"`
	Commence   string `json:"commence_time"`
	Bookmakers []struct {
		Key     string `json:"key"`
		Markets []struct {
			Key      string `json:"key"`
			Outcomes []struct {
				Name  string  `json:"name"`
				Price float64 `json:"price"`
			} `json:"outcomes"`
		} `json:"markets"`
	} `json:"bookmakers"`
}
