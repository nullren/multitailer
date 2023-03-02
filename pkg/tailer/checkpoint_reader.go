package tailer

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type CheckpointReader struct {
	checkpoints map[string]Checkpoint
	sync.Mutex
}

func NewCheckpointReader() *CheckpointReader {
	return &CheckpointReader{
		checkpoints: make(map[string]Checkpoint),
	}
}

const MAX_READ_SIZE = 1024 * 1024 * 1024

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

	reader := io.NewSectionReader(checkpoint.File, checkpoint.Offset, MAX_READ_SIZE)

	scanner, err := NewOffsetScanner(reader)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Printf("%s\t%d\t%d\t%d\n", fileName, len(lines), checkpoint.Offset, scanner.Offset())

	checkpoint.Offset += scanner.Offset()
	r.checkpoints[fileName] = checkpoint
	return lines, nil
}

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
