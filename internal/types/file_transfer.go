package types

import (
	"ShutterSync/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

type FileTransfer struct {
	SourceDrive           Drive
	DestinationDrive      Drive
	Overwrite             bool
	DeleteSourceAfterCopy bool
	OnFileTransferred     func(src, dest string)
	Progress              int
}

func NewFileTransfer(
	source, destination Drive,
	overwrite bool,
	deleteSourceAfterCopy bool,
	onFileTransferred func(src, dest string),

) *FileTransfer {
	return &FileTransfer{
		SourceDrive:           source,
		DestinationDrive:      destination,
		Overwrite:             overwrite,
		DeleteSourceAfterCopy: deleteSourceAfterCopy,
		OnFileTransferred:     onFileTransferred,
		Progress:              0,
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
		if ft.DeleteSourceAfterCopy {
			err = os.Remove(srcFile)
			if err != nil {
				return fmt.Errorf("failed to delete source file %s: %v", srcFile, err)
			}
		}
		ft.Progress += 1
		ft.OnFileTransferred(srcFile, destFile)
	}

	return nil
}
