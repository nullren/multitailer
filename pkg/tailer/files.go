package tailer

import (
	"io/fs"
	"os"
	"sync"
)

type Files struct {
	files []string
	*sync.Mutex
}

func NewFiles(dir string) (Files, error) {
	var files []string
	// TODO: filtering, max depth, etc.
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return Files{files: files}, err
}

func (f *Files) Files() []string {
	f.Lock()
	defer f.Unlock()
	return f.files
}
