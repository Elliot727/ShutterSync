package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// WalkPath traverses the directory at the given path and processes each file or directory.
// It accepts a callback function to handle each file or directory encountered.
func WalkPath(root string, process func(string, os.FileInfo) error) error {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return fmt.Errorf("the path %s does not exist", root)
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return process(path, info)
	})

	return err
}
