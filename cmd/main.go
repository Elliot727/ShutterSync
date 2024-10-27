package main

import (
	"ShutterSync/internal/types"
	"fmt"
	"log"
)

func main() {
	drivePath := "/Volumes/My Passport/Photos"

	// Create a new Drive instance
	drive, err := types.NewDrive(drivePath)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %s", err, drivePath)
	}

	// Provide detailed drive information
	fmt.Printf("Successfully created drive:\n")
	fmt.Printf("Path: %s\n", drive.Path)
	fmt.Printf("Capacity: %d bytes (%.2f GB)\n", drive.Capacity, float64(drive.Capacity)/(1024*1024*1024))
	fmt.Printf("Free Space: %d bytes (%.2f GB)\n", drive.FreeSpace, float64(drive.FreeSpace)/(1024*1024*1024))
	fmt.Printf("Used Space: %d bytes (%.2f GB)\n", drive.UsedSpace, float64(drive.UsedSpace)/(1024*1024*1024))
	fmt.Printf("Is Mounted: %t\n", drive.IsMounted)

	metadataList, err := drive.ExtractMetaData()
	if err != nil {
		log.Fatalf("Error extracting metadata: %v", err)
	}

	// Print the extracted metadata
	fmt.Println("Extracted Metadata:")
	for _, meta := range metadataList {
		fmt.Printf("File: %s, Size: %d bytes, Created: %s, Modified: %s\n",
			meta.FileName, meta.FileSize, meta.CreationTime, meta.ModificationTime)
	}
}
