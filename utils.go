package multitailer

import (
	"context"
	"fmt"
	"time"
)

func PeriodicallyRun(ctx context.Context, interval time.Duration, fn func() error) {
	runner := func() {
		if err := fn(); err != nil {
			fmt.Printf("PeriodicallyRun: failed: %s\n", err)
		}
	}

	timer := time.NewTicker(interval)
	defer timer.Stop()

	// initial run
	runner()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			// periodic run
			runner()
		}
	}
}
