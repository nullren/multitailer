package tailer

import (
	"os"
)

type CheckpointReader struct {
	checkpoints map[string]Checkpoint
}

func NewCheckpointReader() *CheckpointReader {
	return &CheckpointReader{
		checkpoints: make(map[string]Checkpoint),
	}
}

func (r *CheckpointReader) ReadLines(file string) (lines []string, err error) {
	checkpoint := r.checkpoints[file]
	if checkpoint.File == nil {
		checkpoint.File, err = os.Open(file)
		if err != nil {
			return nil, err
		}
	}

	_, err = checkpoint.File.Seek(checkpoint.Offset, 0)
	if err != nil {
		return nil, err
	}

	scanner, err := NewOffsetScanner(checkpoint.File)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	checkpoint.Offset = scanner.Offset()
	return lines, nil
}

func (r *CheckpointReader) Close() error {
	for _, checkpoint := range r.checkpoints {
		if checkpoint.File != nil {
			checkpoint.File.Close()
		}
	}
	return nil
}

type Checkpoint struct {
	Offset int64
	File   *os.File
}
