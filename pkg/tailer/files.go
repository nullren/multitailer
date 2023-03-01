package tailer

import (
	"io/fs"
	"os"
)

type Files struct {
	Files []string `json:"files"`
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
	return Files{files}, err
}
