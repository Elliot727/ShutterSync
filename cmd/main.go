package main

import (
	"ShutterSync/internal/types"
	"log"
)

func main() {
	source := "/Users/elliotsilver/Documents/TESTING 123"
	dest := "/volumes/My Passport"

	// Create a new Drive instance
	sourceDrive, err := types.NewDrive(source)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %v", err, source)
	}

	destDrive, err := types.NewDrive(dest)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %v", err, dest)
	}

	fmt := types.NewFileTransfer(*sourceDrive, *destDrive)

	fmt.Transfer()
}
