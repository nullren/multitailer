package tailer

import (
	"io"
	"os"
	"sync"
)

type CheckpointReader struct {
	checkpoints map[string]Checkpoint

	// maxReadSize is the maximum number of bytes to read from a file.
	// this allows one large file to not dominate the read loop.
	maxReadSize int64
	sync.Mutex
}

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
	if checkpoint.File == nil {
		file, err := os.Open(fileName)
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				return nil, nil
			}
			return nil, err
		}
		checkpoint.File = file
		checkpoint.Offset = 0
	}

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

type Checkpoint struct {
	Offset int64
	File   *os.File
}
