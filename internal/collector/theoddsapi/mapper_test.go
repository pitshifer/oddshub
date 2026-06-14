package theoddsapi

import (
	"testing"
	"time"

	"github.com/pitshifer/oddshub/internal/service"
)

func TestMapToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input []OddsResponse
		want  []service.EventOdds
	}{
		{
			name:  "empty input returns empty result",
			input: []OddsResponse{},
			want:  nil,
		},
		{
			name: "single event without bookmakers",
			input: []OddsResponse{
				{
					ID:       "event-1",
					SportKey: "soccer_epl",
					HomeTeam: "Arsenal",
					AwayTeam: "Chelsea",
					Commence: "2024-01-15T15:00:00Z",
				},
			},
			want: []service.EventOdds{
				{
					EventID:   "event-1",
					Sport:     "soccer_epl",
					HomeTeam:  "Arsenal",
					AwayTeam:  "Chelsea",
					StartTime: time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "event with bookmaker and outcomes",
			input: []OddsResponse{
				{
					ID:       "event-2",
					SportKey: "soccer_epl",
					HomeTeam: "Liverpool",
					AwayTeam: "Man City",
					Commence: "2024-01-20T17:30:00Z",
					Bookmakers: []struct {
						Key     string `json:"key"`
						Markets []struct {
							Key      string `json:"key"`
							Outcomes []struct {
								Name  string  `json:"name"`
								Price float64 `json:"price"`
							} `json:"outcomes"`
						} `json:"markets"`
					}{
						{
							Key: "betway",
							Markets: []struct {
								Key      string `json:"key"`
								Outcomes []struct {
									Name  string  `json:"name"`
									Price float64 `json:"price"`
								} `json:"outcomes"`
							}{
								{
									Key: "h2h",
									Outcomes: []struct {
										Name  string  `json:"name"`
										Price float64 `json:"price"`
									}{
										{Name: "Liverpool", Price: 2.10},
										{Name: "Man City", Price: 3.50},
										{Name: "Draw", Price: 3.20},
									},
								},
							},
						},
					},
				},
			},
			want: []service.EventOdds{
				{
					EventID:   "event-2",
					Sport:     "soccer_epl",
					HomeTeam:  "Liverpool",
					AwayTeam:  "Man City",
					StartTime: time.Date(2024, 1, 20, 17, 30, 0, 0, time.UTC),
					Bookmakers: []service.Bookmaker{
						{
							Name: "betway",
							Markets: []service.Market{
								{
									Type: "h2h",
									Outcomes: []service.Outcome{
										{Name: "Liverpool", Price: 2.10},
										{Name: "Man City", Price: 3.50},
										{Name: "Draw", Price: 3.20},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid commence time falls back to zero time", //FIXME
			input: []OddsResponse{
				{
					ID:       "event-3",
					SportKey: "soccer_epl",
					HomeTeam: "Arsenal",
					AwayTeam: "Chelsea",
					Commence: "не-валидная-дата",
				},
			},
			want: []service.EventOdds{
				{
					EventID:   "event-3",
					Sport:     "soccer_epl",
					HomeTeam:  "Arsenal",
					AwayTeam:  "Chelsea",
					StartTime: time.Time{}, // нулевое время
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOddsToDomain(tt.input)

			// Сравниваем длину результата
			if len(got) != len(tt.want) {
				t.Fatalf("got %d events, want %d", len(got), len(tt.want))
			}

			// Сравниваем каждое событие
			for i := range tt.want {
				assertEventOdds(t, got[i], tt.want[i])
			}
		})
	}
}

// assertEventOdds — вспомогательная функция сравнения.
func assertEventOdds(t *testing.T, got, want service.EventOdds) {
	t.Helper()

	if got.EventID != want.EventID {
		t.Errorf("EventID: got %q, want %q", got.EventID, want.EventID)
	}
	if got.Sport != want.Sport {
		t.Errorf("Sport: got %q, want %q", got.Sport, want.Sport)
	}
	if got.HomeTeam != want.HomeTeam {
		t.Errorf("HomeTeam: got %q, want %q", got.HomeTeam, want.HomeTeam)
	}
	if got.AwayTeam != want.AwayTeam {
		t.Errorf("AwayTeam: got %q, want %q", got.AwayTeam, want.AwayTeam)
	}
	if !got.StartTime.Equal(want.StartTime) {
		t.Errorf("StartTime: got %v, want %v", got.StartTime, want.StartTime)
	}

	if len(got.Bookmakers) != len(want.Bookmakers) {
		t.Fatalf("Bookmakers count: got %d, want %d", len(got.Bookmakers), len(want.Bookmakers))
	}
	for i := range want.Bookmakers {
		assertBookmaker(t, got.Bookmakers[i], want.Bookmakers[i])
	}
}

func assertBookmaker(t *testing.T, got, want service.Bookmaker) {
	t.Helper()

	if got.Name != want.Name {
		t.Errorf("Bookmaker.Name: got %q, want %q", got.Name, want.Name)
	}
	if len(got.Markets) != len(want.Markets) {
		t.Fatalf("Markets count: got %d, want %d", len(got.Markets), len(want.Markets))
	}
	for i := range want.Markets {
		assertMarket(t, got.Markets[i], want.Markets[i])
	}
}

func assertMarket(t *testing.T, got, want service.Market) {
	t.Helper()

	if got.Type != want.Type {
		t.Errorf("Market.Type: got %q, want %q", got.Type, want.Type)
	}
	if len(got.Outcomes) != len(want.Outcomes) {
		t.Fatalf("Outcomes count: got %d, want %d", len(got.Outcomes), len(want.Outcomes))
	}
	for i := range want.Outcomes {
		if got.Outcomes[i] != want.Outcomes[i] {
			t.Errorf("Outcome[%d]: got %+v, want %+v", i, got.Outcomes[i], want.Outcomes[i])
		}
	}
}
