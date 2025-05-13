package utils

import (
	"io"
	"os"
)

const bufferSize = 1 * 1024 * 1024

func CopyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	buf := make([]byte, bufferSize)
	_, err = io.CopyBuffer(output, input, buf)
	if err != nil {
		return err
	}

	return output.Sync()
}
