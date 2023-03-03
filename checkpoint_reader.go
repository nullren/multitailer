package multitailer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// CheckpointReader reads lines from files, starting from the last read position.
// To read a file, call ReadLines with the file name. The file will be opened
// and read from the last read position. The file will be closed when the
// CheckpointReader is closed. If the file is deleted, it will be ignored.
type CheckpointReader struct {
	checkpoints map[string]Checkpoint

	checkpointsSaveFile     string
	checkpointsSaveInterval time.Duration

	// maxReadSize is the maximum number of bytes to read from a file.
	// this allows one large file to not dominate the read loop.
	maxReadSize int64
	sync.Mutex
}

// NewCheckpointReader returns a new CheckpointReader.
func NewCheckpointReader() (*CheckpointReader, error) {
	checkpointsSaveFile := "/tmp/checkpoints.json"
	checkpoints, err := LoadCheckpoints(checkpointsSaveFile)
	if err != nil {
		return nil, fmt.Errorf("load checkpoints failed: %w", err)
	}
	return &CheckpointReader{
		checkpoints:             checkpoints,
		checkpointsSaveFile:     checkpointsSaveFile,
		checkpointsSaveInterval: 30 * time.Second,
		maxReadSize:             1024 * 1024 * 1024,
	}, nil
}

// ReadLines reads lines from a file, starting from the last read position.
func (r *CheckpointReader) ReadLines(fileName string) (lines []string, err error) {
	r.Lock()
	defer r.Unlock()

	checkpoint := r.checkpoints[fileName]
	if err := checkpoint.Check(fileName); err != nil {
		return nil, fmt.Errorf("checkpoint check failed: %w", err)
	}

	// restrict the read size to maxReadSize so that one large file
	// doesn't dominate the read loop.
	reader := io.NewSectionReader(checkpoint.File, checkpoint.Offset, r.maxReadSize)
	scanner, err := NewBytesReadScanner(reader)
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	checkpoint.Offset += scanner.BytesRead()
	r.checkpoints[fileName] = checkpoint
	return lines, nil
}

func (r *CheckpointReader) SaveCheckpoints() error {
	r.Lock()
	defer r.Unlock()
	bytes, err := json.Marshal(r.checkpoints)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	if err := os.WriteFile(r.checkpointsSaveFile, bytes, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}
	fmt.Printf("saved checkpoints: %d bytes written\n", len(bytes))
	return nil
}

func (r *CheckpointReader) RunSaveCheckpoints(ctx context.Context) {
	PeriodicallyRun(ctx, r.checkpointsSaveInterval, r.SaveCheckpoints)
}

// LoadCheckpoints loads checkpoints from a file.
func LoadCheckpoints(checkpoingsFile string) (map[string]Checkpoint, error) {
	bytes, err := os.ReadFile(checkpoingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Checkpoint), nil
		}
		return nil, fmt.Errorf("read file failed: %w", err)
	}
	var checkpoints map[string]Checkpoint
	if err := json.Unmarshal(bytes, &checkpoints); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	fmt.Printf("loaded checkpoints: %d bytes read\n", len(bytes))
	return checkpoints, nil
}

// Close closes all open files.
func (r *CheckpointReader) Close() error {
	r.Lock()
	defer r.Unlock()

	for _, checkpoint := range r.checkpoints {
		if checkpoint.File != nil {
			checkpoint.File.Close()
			checkpoint.File = nil
		}
	}
	return nil
}
