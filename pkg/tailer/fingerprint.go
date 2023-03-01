package tailer

import (
	"bufio"
	"hash/crc32"
	"io"
	"os"
)

// Fingerprint is a struct that contains a file name and a hash of the first line of the file.
type Fingerprint struct {
	File          string `json:"file"`
	FirstLineHash uint32 `json:"first_line_hash"`
}

const MAX_LINE_LENGTH = 1024

func NewFingerprint(fileName string) (fingerPrint Fingerprint, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return fingerPrint, err
	}
	defer file.Close()

	// get first line or MAX_LINE_LENGTH bytes
	reader := io.LimitReader(file, MAX_LINE_LENGTH)
	scanner := bufio.NewScanner(reader)

	if scanner.Scan() {
		buffer := scanner.Bytes()
		fingerPrint = Fingerprint{
			File:          fileName,
			FirstLineHash: crc32.Checksum(buffer, crc32.MakeTable(crc32.Castagnoli)),
		}
		return fingerPrint, nil
	}

	// empty file if Err nil
	return fingerPrint, scanner.Err()
}
