package types

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Drive represents a storage device used for file transfers.

type Drive struct {
	Path      string
	Capacity  uint64
	FreeSpace uint64
	UsedSpace uint64
	IsMounted bool
}

// NewDrive initializes a new Drive instance based on the provided path.
func NewDrive(path string) (*Drive, error) {
	d := Drive{}
	drive, err := d.GetDriveDetails(path)
	if err != nil {
		return nil, err
	}
	return &drive, nil
}

// GetDriveDetails retrieves the filesystem statistics for the given path.

func (d *Drive) GetDriveDetails(path string) (Drive, error) {
	var stat syscall.Statfs_t

	// Get filesystem statistics
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return Drive{}, err
	}

	// Capacity is total space (block size * total blocks)
	capacity := uint64(stat.Bsize) * stat.Blocks
	freeSpace := uint64(stat.Bsize) * stat.Bavail
	usedSpace := capacity - freeSpace

	drive := Drive{
		Path:      path,
		Capacity:  capacity,
		FreeSpace: freeSpace,
		UsedSpace: usedSpace,
		IsMounted: true,
	}

	return drive, nil
}

// WalkPath traverses the directory at the given path and processes each file or directory.
// It accepts a callback function to handle each file or directory encountered.
func (d *Drive) WalkPath(process func(string, os.FileInfo) error) error {
	if _, err := os.Stat(d.Path); os.IsNotExist(err) {
		return fmt.Errorf("the path %s does not exist", d.Path)
	}

	err := filepath.Walk(d.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return process(path, info)
	})

	return err
}

// Metadata holds the extracted metadata from a file.

// ExtractMetaData extracts metadata from files in the Drive's directory.
func (d *Drive) ExtractMetaData() ([]Metadata, error) {
	var metadataList []Metadata

	// Walk the directory and extract metadata from each file
	err := filepath.Walk(d.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a Metadata instance for the file
		meta := Metadata{
			FileName:         info.Name(),
			FileSize:         uint64(info.Size()),
			CreationTime:     info.ModTime().Format("2006-01-02 15:04:05"), // Customize as needed
			ModificationTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}

		// Append to the metadata list
		metadataList = append(metadataList, meta)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk the directory: %v", err)
	}

	return metadataList, nil
}
