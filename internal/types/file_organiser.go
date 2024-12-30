package types

import (
	"fmt"
	"os"
	"path/filepath"
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
	// Get a list of files in the drive directory
	files, err := os.ReadDir(fo.Drive.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			// Skip directories, or handle them as needed
			continue
		}
		// Ensure the file name is long enough to extract a date (length 10 for 'DD-MM-YYYY')
		if len(file.Name()) < 10 {
			fmt.Printf("Skipping file %s: Filename too short to extract date\n", file.Name())
			continue // Skip this file if it's too short
		}
		// Extract the date from the filename
		dateStr := file.Name()[:10]
		date, err := time.Parse("02-01-2006", dateStr)
		if err != nil {
			return fmt.Errorf("failed to parse date for file %s: %v", file.Name(), err)
		}

		// Create directories for year, month, and day
		yearDir := filepath.Join(fo.Drive.Path, fmt.Sprintf("%d", date.Year()))
		monthDir := filepath.Join(yearDir, fmt.Sprintf("%02d", date.Month()))
		dayDir := filepath.Join(monthDir, fmt.Sprintf("%02d", date.Day()))

		// Ensure the directory structure exists
		err = os.MkdirAll(dayDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directories for %s: %v", file.Name(), err)
		}

		// Move the file to the organized directory
		srcFile := filepath.Join(fo.Drive.Path, file.Name())
		destFile := filepath.Join(dayDir, file.Name())
		err = os.Rename(srcFile, destFile)
		if err != nil {
			return fmt.Errorf("failed to move %s to %s: %v", srcFile, destFile, err)
		}

		// Callback for each organized file
		fo.OnFileOrganized(srcFile, destFile)
	}

	return nil
}
