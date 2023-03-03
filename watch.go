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

// Watch files for new content and call watchFunc for each new line of content.
// Files are searched for using searchGlob.
func Watch(ctx context.Context, searchGlob string, watchFunc WatchFunc) error {
	files := NewFiles(searchGlob)
	go files.RunUpdater(ctx)

	reader, err := NewCheckpointReader()
	if err != nil {
		return err
	}
	defer reader.Close()
	go reader.RunSaveCheckpoints(ctx)

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
