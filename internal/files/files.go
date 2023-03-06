package files

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/nullren/multitailer/internal/utils"
)

// Store is a struct that holdes the current files that should be tailed.
type Store struct {
	searchGlob          string
	filesUpdateInterval time.Duration
	files               []string
	sync.Mutex
}

// NewStore creates a new files Store that will search for files matching the
// searchGlob and refresh the list of files every filesUpdateInterval.
func NewStore(searchGlob string, filesUpdateInterval time.Duration) *Store {
	return &Store{
		searchGlob:          searchGlob,
		filesUpdateInterval: filesUpdateInterval,
	}
}

// UpdateFiles updates the list of files to tail.
func (f *Store) UpdateFiles() error {
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

// RunUpdater runs the UpdateFiles function periodically.
func (f *Store) RunUpdater(ctx context.Context) {
	utils.PeriodicallyRun(ctx, f.filesUpdateInterval, f.UpdateFiles)
}

// Files returns the list of files to tail.
func (f *Store) Files() []string {
	f.Lock()
	defer f.Unlock()

	return f.files
}
