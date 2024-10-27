package main

import (
	"ShutterSync/internal/types"
	"fmt"
	"log"
)

func main() {
	source := "/Users/elliotsilver/Documents/TESTING 123"
	dest := "/volumes/My Passport"

	sourceDrive, err := types.NewDrive(source)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %v", err, source)
	}

	destDrive, err := types.NewDrive(dest)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %v", err, dest)
	}

	fmt := types.NewFileTransfer(
		*sourceDrive,
		*destDrive,
		true,
		false,
		func(src, dest string) {
			fmt.Printf("Transferred %s to %s\n", src, dest)
		},
		func(file string, err error) {
			fmt.Printf("Error transferring %s: %v\n", file, err)
		},
	)

	fmt.Transfer()
}
