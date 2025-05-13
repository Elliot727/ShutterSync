package types

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileOrganizer struct {
	Drive           Drive
	OnFileOrganized func(src, dest string)
}

func NewFileOrganizer(drive Drive, onFileOrganized func(src, dest string)) *FileOrganizer {
	return &FileOrganizer{
		Drive:           drive,
		OnFileOrganized: onFileOrganized,
	}
}

func (fo *FileOrganizer) Organize() error {

	files, err := os.ReadDir(fo.Drive.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	numWorkers := runtime.NumCPU() * 4
	fileChan := make(chan string, numWorkers*2)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var dirCache sync.Map
	var errors []error

	appendError := func(e error) {
		mu.Lock()
		errors = append(errors, e)
		mu.Unlock()
	}

	worker := func() {
		defer wg.Done()
		for fileName := range fileChan {
			srcFile := filepath.Join(fo.Drive.Path, fileName)

			if len(fileName) < 10 {
				fmt.Printf("Skipping file %s: filename too short to contain a date\n", fileName)
				continue
			}

			dateStr := fileName[:10]
			date, err := time.Parse("02-01-2006", dateStr)
			if err != nil {
				fmt.Printf("Skipping file %s: invalid date format: %v\n", fileName, err)
				continue
			}

			dayDir := filepath.Join(
				fo.Drive.Path,
				fmt.Sprintf("%d/%02d/%02d", date.Year(), date.Month(), date.Day()),
			)

			if _, exists := dirCache.Load(dayDir); !exists {
				if err := os.MkdirAll(dayDir, os.ModePerm); err != nil {
					appendError(fmt.Errorf("failed to create directory %s: %w", dayDir, err))
					continue
				}
				dirCache.Store(dayDir, struct{}{})
			}

			destFile := filepath.Join(dayDir, fileName)
			if err := os.Rename(srcFile, destFile); err != nil {
				appendError(fmt.Errorf("failed to move %s to %s: %w", srcFile, destFile, err))
				continue
			}

			if fo.OnFileOrganized != nil {
				fo.OnFileOrganized(srcFile, destFile)
			}
		}
	}

	for range numWorkers {
		wg.Add(1)
		go worker()
	}

	for _, file := range files {
		if !file.IsDir() {
			fileChan <- file.Name()
		}
	}
	close(fileChan)

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("organization completed with %d errors", len(errors))
	}

	return nil
}
