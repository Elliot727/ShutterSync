package types

import (
	"ShutterSync/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type FileTransfer struct {
	SourceDrive           Drive
	DestinationDrive      Drive
	DeleteSourceAfterCopy bool
	OnFileTransferred     func(src, dest string)
}

type fileJob struct {
	name string
	info os.FileInfo
}

func NewFileTransfer(source, destination Drive, deleteSourceAfterCopy bool, onFileTransferred func(src, dest string)) *FileTransfer {
	return &FileTransfer{
		SourceDrive:           source,
		DestinationDrive:      destination,
		DeleteSourceAfterCopy: deleteSourceAfterCopy,
		OnFileTransferred:     onFileTransferred,
	}
}

func (ft *FileTransfer) Transfer() error {
	files, err := os.ReadDir(ft.SourceDrive.Path)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	numWorkers := runtime.NumCPU() * 4
	fileChan := make(chan fileJob, numWorkers*2)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	appendError := func(e error) {
		mu.Lock()
		errors = append(errors, e)
		mu.Unlock()
	}

	worker := func() {
		defer wg.Done()
		for job := range fileChan {
			fileName := job.name
			info := job.info

			srcFile := filepath.Join(ft.SourceDrive.Path, fileName)
			metadata, err := ft.SourceDrive.ExtractMetaData(info)
			if err != nil {
				appendError(fmt.Errorf("skipping %s: %w", fileName, err))
				continue
			}

			ext := filepath.Ext(fileName)
			baseName := metadata.CreationTime
			destDir := ft.DestinationDrive.Path
			destFile := filepath.Join(destDir, baseName+ext)
			counter := 0

			var output *os.File
			for {
				output, err = os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
				if err == nil {
					break
				}
				if os.IsExist(err) {
					counter++
					destFile = filepath.Join(destDir, fmt.Sprintf("%s%d%s", baseName, counter, ext))
				} else {
					appendError(fmt.Errorf("error creating file %s: %w", destFile, err))
					break
				}
			}
			if output == nil {
				continue
			}
			output.Close()

			if err := utils.CopyFile(srcFile, destFile); err != nil {
				appendError(fmt.Errorf("error copying %s: %w", srcFile, err))
				os.Remove(destFile)
				continue
			}

			if ft.DeleteSourceAfterCopy {
				if err := os.Remove(srcFile); err != nil {
					appendError(fmt.Errorf("error deleting %s: %w", srcFile, err))
				}
			}

			if ft.OnFileTransferred != nil {
				ft.OnFileTransferred(srcFile, destFile)
			}
		}
	}

	for range numWorkers {
		wg.Add(1)
		go worker()
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			appendError(fmt.Errorf("failed to get info for %s: %w", file.Name(), err))
			continue
		}
		fileChan <- fileJob{name: file.Name(), info: info}
	}

	close(fileChan)
	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("transfer completed with %d errors", len(errors))
	}
	return nil
}
