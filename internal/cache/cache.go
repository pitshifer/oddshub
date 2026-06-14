package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/pitshifer/oddshub/internal/service"
)

type Cache struct {
	mu      sync.RWMutex
	storage service.Storage

	bySport map[string][]service.EventOdds
	prices  map[string]float64
}

func New(storage service.Storage) *Cache {
	return &Cache{
		storage: storage,
		bySport: make(map[string][]service.EventOdds),
		prices:  make(map[string]float64),
	}
}

func (c *Cache) GetOdds(ctx context.Context, sport string) ([]service.EventOdds, error) {
	c.mu.RLock()
	odds, ok := c.bySport[sport]
	c.mu.RUnlock()
	if ok {
		return odds, nil
	}

	// get from storage
	odds, err := c.storage.GetOdds(ctx, sport)
	if err != nil {
		return nil, err
	}

	// save to cache
	c.mu.Lock()
	c.bySport[sport] = odds
	c.mu.Unlock()

	return odds, nil
}

func (c *Cache) SaveOdds(ctx context.Context, provider string, odds []service.EventOdds) error {
	c.mu.RLock()
	changeSet := make(map[string]service.EventOdds)
	for _, event := range odds {
		for _, bm := range event.Bookmakers {
			for _, market := range bm.Markets {
				for _, outcome := range market.Outcomes {
					key := fmt.Sprintf("%s:%s:%s:%s", event.EventID, bm.Name, market.Type, outcome.Name)
					if cached, ok := c.prices[key]; !ok || cached != outcome.Price {
						changeSet[event.EventID] = event
					}
				}
			}
		}
	}
	c.mu.RUnlock()

	if len(changeSet) == 0 {
		return nil
	}

	// save to storage (DB)
	changes := make([]service.EventOdds, 0, len(changeSet))
	for _, event := range changeSet {
		changes = append(changes, event)
	}
	if err := c.storage.SaveOdds(ctx, provider, changes); err != nil {
		return err
	}

	c.mu.Lock()
	for _, event := range changes {
		c.bySport[event.Sport] = updateEvents(c.bySport[event.Sport], event)
		for _, bm := range event.Bookmakers {
			for _, market := range bm.Markets {
				for _, outcome := range market.Outcomes {
					key := fmt.Sprintf("%s:%s:%s:%s", event.EventID, bm.Name, market.Type, outcome.Name)
					c.prices[key] = outcome.Price
				}
			}
		}
	}
	c.mu.Unlock()

	return nil
}

func (c *Cache) SaveSports(ctx context.Context, sports []service.Sport) error {
	return c.storage.SaveSports(ctx, sports)
}

func (c *Cache) GetSports(ctx context.Context) ([]service.Sport, error) {
	return c.storage.GetSports(ctx)
}

func updateEvents(events []service.EventOdds, updated service.EventOdds) []service.EventOdds {
	for i, e := range events {
		if e.EventID == updated.EventID {
			events[i] = updated
			return events
		}
	}
	return append(events, updated)
}

func (c *Cache) Warm(ctx context.Context, sports []string) error {
	bySport := make(map[string][]service.EventOdds)
	prices := make(map[string]float64)

	for _, sport := range sports {
		odds, err := c.storage.GetOdds(ctx, sport)
		if err != nil {
			return err
		}

		// save to cache
		bySport[sport] = odds

		// save prices to compare later
		for _, event := range odds {
			for _, bm := range event.Bookmakers {
				for _, market := range bm.Markets {
					for _, outcome := range market.Outcomes {
						key := fmt.Sprintf("%s:%s:%s:%s", event.EventID, bm.Name, market.Type, outcome.Name)
						prices[key] = outcome.Price
					}
				}
			}
		}
	}

	c.mu.Lock()
	c.bySport = bySport
	c.prices = prices
	c.mu.Unlock()

	return nil
}
