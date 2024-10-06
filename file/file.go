package file

import (
	"os"
)

func CreateDirIfNotExists(path string) error {
	return os.MkdirAll(path, 0755)
}
