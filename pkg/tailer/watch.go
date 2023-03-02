package tailer

import (
	"context"
	"fmt"
	"time"
)

// WatchFunc is a function that is called when a file has
// new content. It is called for each new line data.
// Parameters are (file, line). Returning an error will
// end the watch.
type WatchFunc = func(string, string) error

func Watch(ctx context.Context, dir string, watchFunc WatchFunc) error {
	files, err := NewFiles(dir)
	if err != nil {
		return err
	}

	reader := NewCheckpointReader()

	for {
		for _, file := range files.Files() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fmt.Println("reading", file)
				lines, err := reader.ReadLines(file)
				if err != nil {
					return err
				}
				for _, line := range lines {
					if err := watchFunc(file, line); err != nil {
						return err
					}
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
