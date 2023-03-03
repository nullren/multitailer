package multitailer

import (
	"fmt"
	"os"
	"syscall"
)

type Checkpoint struct {
	Offset int64
	Size   int64
	Dev    int32
	Ino    uint64
	File   *os.File `json:"-"`
}

// Check checks if the file has changed, and re-opens it if it has.
func (c *Checkpoint) Check(fileName string) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			// file removed or can't read it anymore
			_ = c.File.Close()
			return nil
		}
		return fmt.Errorf("stat failed: %w", err)
	}
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("not a syscall.Stat_t: %v", fileInfo.Sys())
	}

	// if we loaded a checkpoint from a file, the file may not be open yet
	if c.File == nil && c.Offset > 0 {
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("open failed: %w", err)
		}
		if _, err := file.Seek(c.Offset, 0); err != nil {
			return fmt.Errorf("seek failed: %w", err)
		}
		c.File = file
	}

	// if checkpoint is new (all values 0), or file was moved or truncated, re-open it

	if c.Size > fileInfo.Size() || c.Dev != stat.Dev || c.Ino != stat.Ino {
		if c.File != nil {
			_ = c.File.Close()
		}
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("open failed: %w", err)
		}
		c.File = file
		c.Offset = 0
		c.Dev = stat.Dev
		c.Ino = stat.Ino
	}
	c.Size = fileInfo.Size()
	return nil
}
