package multitailer

import (
	"context"
	"fmt"
	"time"
)

type WatchConfig struct {
	// CheckpointsSaveFile is the file to save file read checkpoints to.
	CheckpointsSaveFile string

	// CheckpointsSaveInterval is the interval to save checkpoints to disk.
	CheckpointsSaveInterval time.Duration

	// FileSearchGlob is the glob to search for files to watch.
	FileSearchGlob string

	// FileUpdateInterval is the interval to update the list of files to watch.
	FileUpdateInterval time.Duration

	// FileMaxReadBytes is the maximum number of bytes to read from a file
	FileMaxReadBytes int64
}

type Watch struct {
	files  *Files
	reader *CheckpointReader
}

func NewWatch(config WatchConfig) (*Watch, error) {
	files := NewFiles(config.FileSearchGlob, config.FileUpdateInterval)
	reader, err := NewCheckpointReader(CheckpointConfig{
		SaveFile:     config.CheckpointsSaveFile,
		SaveInterval: config.CheckpointsSaveInterval,
		MaxReadBytes: config.FileMaxReadBytes,
	})
	if err != nil {
		return nil, err
	}
	return &Watch{
		files:  files,
		reader: reader,
	}, nil
}

// WatchFunc is a function that is called when a file has
// new content. It is called for each new line data.
// Parameters are (file, line). Returning an error will
// end the watch.
type WatchFunc = func(file, line string) error

// Watch files for new content and call watchFunc for each new line of content.
// Files are searched for using searchGlob.
func (w *Watch) Watch(ctx context.Context, watchFunc WatchFunc) error {
	go w.files.RunUpdater(ctx)
	defer w.reader.Close()
	go w.reader.RunSaveCheckpoints(ctx)

	for {
		for _, file := range w.files.Files() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				lines, err := w.reader.ReadLines(file)
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
