package tailer

import (
	"bufio"
	"bytes"
	"io"
)

// BytesReadScanner is a bufio.Scanner that also keeps track of the number of bytes read.
type BytesReadScanner struct {
	bytesRead int64
	*bufio.Scanner
}

// NewBytesReadScanner returns a new BytesReadScanner that reads from r.
func NewBytesReadScanner(r io.Reader) (*BytesReadScanner, error) {
	scanner := bufio.NewScanner(r)
	bytesReadScanner := &BytesReadScanner{
		Scanner: scanner,
	}

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		defer func() {
			bytesReadScanner.bytesRead += int64(advance)
		}()
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			return i + 1, dropCR(data[0:i]), nil
		}
		if atEOF {
			// If we're at EOF, we have a final, non-terminated line.
			// We should leave it alone until more data comes in.
			return 0, nil, nil
		}
		return 0, nil, nil
	})

	return bytesReadScanner, nil
}

// BytesRead returns the number of bytes read by the scanner.
func (s *BytesReadScanner) BytesRead() int64 {
	return s.bytesRead
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
