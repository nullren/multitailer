package tailer

import (
	"bufio"
	"bytes"
	"io"
)

type OffsetScanner struct {
	offset int64
	*bufio.Scanner
}

func NewOffsetScanner(r io.Reader) (*OffsetScanner, error) {
	scanner := bufio.NewScanner(r)
	offsetScanner := &OffsetScanner{
		Scanner: scanner,
	}

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		defer func() {
			offsetScanner.offset += int64(advance)
		}()
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			return i + 1, dropCR(data[0:i]), nil
		}
		if atEOF {
			return len(data), dropCR(data), nil
		}
		return 0, nil, nil
	})

	return offsetScanner, nil
}

func (s *OffsetScanner) Offset() int64 {
	return s.offset
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
