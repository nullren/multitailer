package checkpointer

import (
	"fmt"
	"os"
	"syscall"
)

type Checkpoint struct {
	Offset int64
	Size   int64
	Dev    uint64
	Ino    uint64
	File   *os.File `json:"-"`
}

// Check checks if the file has changed, and re-opens it if it has.
func (c *Checkpoint) Check(fileName string) error {
	// if the file is a symlink, stat will return the info for the target the
	// symlink points to — which is what we want. Lstat would return the info
	// for the symlink itself.
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

	// we loaded a checkpoint from a file, so the file is not open yet
	if c.File == nil && c.Offset > 0 {
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("open failed: %w", err)
		}
		if _, err := file.Seek(c.Offset, 0); err != nil {
			return fmt.Errorf("seek failed: %w", err)
		}
		c.File = file
		// setting these avoids re-opening in the next check
		c.Dev = uint64(stat.Dev)
		c.Ino = stat.Ino
		c.Size = fileInfo.Size()
	}

	// file was moved or truncated, re-open it
	if c.Size > fileInfo.Size() || c.Dev != uint64(stat.Dev) || c.Ino != stat.Ino {
		if c.File != nil {
			_ = c.File.Close()
		}
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("open failed: %w", err)
		}
		c.File = file
		c.Offset = 0
		c.Dev = uint64(stat.Dev)
		c.Ino = stat.Ino
	}

	// always update the size
	c.Size = fileInfo.Size()
	return nil
}
