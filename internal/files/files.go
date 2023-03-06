package files

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/nullren/multitailer/internal/utils"
)

type Files struct {
	searchGlob          string
	filesUpdateInterval time.Duration
	files               []string
	sync.Mutex
}

func NewFiles(searchGlob string, filesUpdateInterval time.Duration) *Files {
	return &Files{
		searchGlob:          searchGlob,
		filesUpdateInterval: filesUpdateInterval,
	}
}

func (f *Files) UpdateFiles() error {
	f.Lock()
	defer f.Unlock()

	files, err := filepath.Glob(f.searchGlob)
	if err != nil {
		return err
	}

	fmt.Printf("found files: %v\n", files)
	f.files = files
	return nil
}

func (f *Files) RunUpdater(ctx context.Context) {
	utils.PeriodicallyRun(ctx, f.filesUpdateInterval, f.UpdateFiles)
}

func (f *Files) Files() []string {
	f.Lock()
	defer f.Unlock()

	return f.files
}
