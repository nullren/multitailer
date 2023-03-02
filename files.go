package multitailer

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

type Files struct {
	searchGlob          string
	filesUpdateInterval time.Duration
	files               []string
	sync.Mutex
}

func NewFiles(searchGlob string) *Files {
	return &Files{
		searchGlob:          searchGlob,
		filesUpdateInterval: 10 * time.Second,
	}
}

func (f *Files) UpdateFiles() error {
	f.Lock()
	defer f.Unlock()

	files, err := filepath.Glob(f.searchGlob)
	if err != nil {
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
