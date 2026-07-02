package service

import (
	"context"
	"time"
)

type Storage interface {
	SaveSports(ctx context.Context, sports []Sport) error
	GetSports(ctx context.Context) ([]Sport, error)

	SaveOdds(ctx context.Context, provider string, odds []EventOdds) error
	GetOdds(ctx context.Context, sport string) ([]EventOdds, error)
}

type Odds interface {
	GetSports(ctx context.Context) ([]Sport, error)
	GetOdds(ctx context.Context, sport string) ([]EventOdds, error)
}

type Sport struct {
	Key          string `json:"key"`
	Title        string `json:"title"`
	Group        string `json:"group"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
	HasOutrights bool   `json:"hasOutrights"`
}

type Leagues struct {
	ID   int
	Name string
}

type EventOdds struct {
	EventID    string      `json:"eventId"`
	Sport      string      `json:"sport"`
	HomeTeam   string      `json:"homeTeam"`
	AwayTeam   string      `json:"awayTeam"`
	StartTime  time.Time   `json:"startTime"`
	Bookmakers []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Name    string   `json:"name"`
	Markets []Market `json:"markets"`
}

type Market struct {
	Type     string    `json:"type"`
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
