package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDirIfNotExists(path string) error {
	return os.MkdirAll(path, 0755)
}

func WriteFileWithDir(filename string, contents []byte) error {
	if err := CreateDirIfNotExists(filepath.Dir(filename)); err != nil {
		return fmt.Errorf("failed to create directory for file %q: %w", filename, err)
	}
	return os.WriteFile(filename, contents, 0644)
}
