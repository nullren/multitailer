package tailer

type CheckpointReader struct{}

func NewCheckpointReader() *CheckpointReader {
	return &CheckpointReader{}
}

func (r *CheckpointReader) ReadLines(file string) (lines []string, err error) {
	return nil, nil
}
