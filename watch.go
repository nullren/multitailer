package multitailer

import (
	"context"
	"fmt"
	"time"
)

// WatchFunc is a function that is called when a file has
// new content. It is called for each new line data.
// Parameters are (file, line). Returning an error will
// end the watch.
type WatchFunc = func(file, line string) error

// Watch watches a directory for files and new content and calls
// watchFunc for each new line of content.
func Watch(ctx context.Context, dir string, watchFunc WatchFunc) error {
	files := NewFiles(dir)
	go files.RunUpdater(ctx)

	reader := NewCheckpointReader()
	defer reader.Close()

	for {
		for _, file := range files.Files() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				lines, err := reader.ReadLines(file)
				if err != nil {
					return err
				}
				for _, line := range lines {
					if err := watchFunc(file, line); err != nil {
						return err
					}
				}
				fmt.Printf("read %d lines from %s\n", len(lines), file)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
