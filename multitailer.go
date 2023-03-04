package multitailer

import (
	"context"
	"fmt"
	"time"
)

type MultitailerConfig struct {
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

// FollowFunc is a function that is called when a file has new content. It is
// called for each new line data. Parameters are (file, line). Returning an
// error will end the outer Follow function calling this.
type FollowFunc = func(file, line string) error

type Multitailer struct {
	files  *Files
	reader *CheckpointReader
}

func NewMultitailer(config MultitailerConfig) (*Multitailer, error) {
	files := NewFiles(config.FileSearchGlob, config.FileUpdateInterval)
	reader, err := NewCheckpointReader(CheckpointConfig{
		SaveFile:     config.CheckpointsSaveFile,
		SaveInterval: config.CheckpointsSaveInterval,
		MaxReadBytes: config.FileMaxReadBytes,
	})
	if err != nil {
		return nil, err
	}
	return &Multitailer{
		files:  files,
		reader: reader,
	}, nil
}

// Follow watches files for new content and call followFunc for each new line of
// content. Files are searched for using searchGlob.
func (m *Multitailer) Follow(ctx context.Context, followFunc FollowFunc) error {
	go m.files.RunUpdater(ctx)
	defer m.reader.Close()
	go m.reader.RunSaveCheckpoints(ctx)

	for {
		for _, file := range m.files.Files() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				lines, err := m.reader.ReadLines(file)
				if err != nil {
					return err
				}
				for _, line := range lines {
					if err := followFunc(file, line); err != nil {
						return err
					}
				}
				fmt.Printf("read %d lines from %s\n", len(lines), file)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
