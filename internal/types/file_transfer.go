package types

import (
	"ShutterSync/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

type FileTransfer struct {
	SourceDrive      Drive
	DestinationDrive Drive
}

func NewFileTransfer(source, destination Drive) *FileTransfer {
	return &FileTransfer{
		SourceDrive:      source,
		DestinationDrive: destination,
	}
}

func (ft *FileTransfer) Transfer() error {
	// Get a list of files in the source directory
	files, err := os.ReadDir(ft.SourceDrive.Path)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			// Skip directories, or handle them as needed
			continue
		}

		// Get the full source file path
		srcFile := filepath.Join(ft.SourceDrive.Path, file.Name())

		// Extract metadata (assuming you have this function implemented)
		metadata, err := ft.SourceDrive.ExtractMetaData(srcFile)
		if err != nil {
			return fmt.Errorf("failed to extract metadata for %s: %v", file.Name(), err)
		}
		ext := filepath.Ext(file.Name())

		// Generate a new file name based on metadata (e.g., using timestamp)

		newFileName := fmt.Sprintf("%s%s", metadata.CreationTime, ext)

		destFile := filepath.Join(ft.DestinationDrive.Path, newFileName)
		fmt.Println(destFile, srcFile)
		// Copy the file to the destination
		err = utils.CopyFile(srcFile, destFile)
		if err != nil {
			return fmt.Errorf("failed to copy %s to %s: %v", srcFile, destFile, err)
		}
	}

	return nil
}
