package tailer

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Files struct {
	files []string
}

func NewFiles(dir string) (*Files, error) {
	var files []string
	// TODO: filtering, max depth, etc.
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fullPath := filepath.Join(dir, path)
			files = append(files, fullPath)
		}
		return nil
	})
	return &Files{files: files}, err
}

func (f *Files) Files() []string {
	return f.files
}
