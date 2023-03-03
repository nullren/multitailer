package multitailer

import (
	"context"
	"fmt"
	"time"
)

func PeriodicallyRun(ctx context.Context, interval time.Duration, fn func() error) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	// initial run
	if err := fn(); err != nil {
		fmt.Printf("PeriodicallyRun failed: %s\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := fn(); err != nil {
				fmt.Printf("PeriodicallyRun: failed: %s\n", err)
			}
		}
	}
}
