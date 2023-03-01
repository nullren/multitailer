package tailer

import (
	"context"
)

// WatchFunc is a function that is called when a file has
// new content. Parameters are (file, content).
type WatchFunc = func(string, string) error

func Watch(ctx context.Context, dir string, watchFunc WatchFunc) error {
	files, err := NewFiles(dir)
	if err != nil {
		return err
	}

	for {
		for _, file := range files.Files {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// TODO: get next line of file
				// lines := file.ReadLines()
				if err := watchFunc(file, ""); err != nil {
					return err
				}
			}
			// time.Sleep(1 * time.Second)
		}
	}
}
