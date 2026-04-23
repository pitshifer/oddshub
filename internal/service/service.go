package service

import (
	"context"
	"time"
)

type Storage interface {
	SaveOdds(ctx context.Context, provider string, odds []EventOdds) error
	GetOdds(ctx context.Context, sport string) ([]EventOdds, error)
}

type Odds interface {
	// GetSports(ctx context.Context) ([]Sport, error)
	// GetLeagues(ctx context.Context, sport string) ([]Leagues, error)
	GetOdds(ctx context.Context, sport string) ([]EventOdds, error)
}

type Sport struct {
	Key  string
	Name string
}

type Leagues struct {
	ID   int
	Name string
}

type EventOdds struct {
	EventID    string
	Sport      string
	HomeTeam   string
	AwayTeam   string
	StartTime  time.Time
	Bookmakers []Bookmaker
}

type Bookmaker struct {
	Name    string
	Markets []Market
}

type Market struct {
	Type     string
	Outcomes []Outcome
}

type Outcome struct {
	Name  string
	Price float64
}
