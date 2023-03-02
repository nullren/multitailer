package tailer

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Files struct {
	searchDir           string
	filesUpdateInterval time.Duration
	files               []string
	sync.Mutex
}

func NewFiles(dir string) *Files {
	return &Files{
		searchDir:           dir,
		filesUpdateInterval: 10 * time.Second,
	}
}

func (f *Files) UpdateFiles() error {
	f.Lock()
	defer f.Unlock()

	var files []string
	// TODO: filtering, max depth, etc.
	if err := fs.WalkDir(os.DirFS(f.searchDir), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fullPath := filepath.Join(f.searchDir, path)
			files = append(files, fullPath)
		}
		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("found files: %v", files)
	f.files = files
	return nil
}

func (f *Files) RunUpdater(ctx context.Context) {
	timer := time.NewTicker(f.filesUpdateInterval)
	defer timer.Stop()

	// initial update
	_ = f.UpdateFiles()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			_ = f.UpdateFiles()
		}
	}
}

func (f *Files) Files() []string {
	f.Lock()
	defer f.Unlock()

	return f.files
}
