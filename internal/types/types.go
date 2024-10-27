package types

import "syscall"

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
