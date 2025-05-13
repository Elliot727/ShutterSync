package types

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Drive struct {
	Path      string
	Capacity  uint64
	FreeSpace uint64
	UsedSpace uint64
	IsMounted bool
}

func NewDrive(path string) (*Drive, error) {
	d := Drive{}
	drive, err := d.GetDriveDetails(path)
	if err != nil {
		return nil, err
	}
	return &drive, nil
}

func (d *Drive) GetDriveDetails(path string) (Drive, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(path, &stat)
	if err != nil {
		return Drive{}, err
	}

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

func (d *Drive) ExtractMetaData(info os.FileInfo) (Metadata, error) {
	if info.IsDir() {
		return Metadata{}, fmt.Errorf("path is a directory: %s", info.Name())
	}

	return Metadata{
		FileName:         info.Name(),
		FileSize:         uint64(info.Size()),
		CreationTime:     info.ModTime().Format("02-01-2006_15-04"),
		ModificationTime: info.ModTime().Format("02-01-2006_15-04"),
	}, nil
}
