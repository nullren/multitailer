package multitailer

import (
	"fmt"
	"io"
	"sync"
)

// CheckpointReader reads lines from files, starting from the last read position.
// To read a file, call ReadLines with the file name. The file will be opened
// and read from the last read position. The file will be closed when the
// CheckpointReader is closed. If the file is deleted, it will be ignored.
type CheckpointReader struct {
	checkpoints map[string]Checkpoint

	// maxReadSize is the maximum number of bytes to read from a file.
	// this allows one large file to not dominate the read loop.
	maxReadSize int64
	sync.Mutex
}

// NewCheckpointReader returns a new CheckpointReader.
func NewCheckpointReader() *CheckpointReader {
	return &CheckpointReader{
		checkpoints: make(map[string]Checkpoint),
		maxReadSize: 1024 * 1024 * 1024,
	}
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
