package users

import (
	"context"
	"log"
	"time"
)

// Ticker ...
type Ticker struct {
	cleanuper *Cleanuper
	interval  time.Duration
}

// NewTicker ...
func NewTicker(cleanuper *Cleanuper, interval time.Duration) *Ticker {
	return &Ticker{
		cleanuper: cleanuper,
		interval:  interval,
	}
}

// Run ...
func (t *Ticker) Run(ctx context.Context) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[users_ticker] context is done. Exit ticker...")
			return ctx.Err()
		case <-ticker.C:
			log.Printf("[users_ticker] cleanup users step...")
			err := t.cleanuper.Cleanup(ctx)
			if err != nil {
				log.Printf("[users_ticker] users cleanup finished with error: %v", err)
			} else {
				log.Printf("[users_ticker] users cleanup finished successfully")
			}
		}
	}
}
