package main

import (
	"ShutterSync/internal/types"
	"fmt"
	"log"
	"os"

	"github.com/sqweek/dialog"
)

func main() {

	// Show file picker for source directory
	srcDir, err := dialog.Directory().Title("Select Source Directory").Browse()
	if err != nil {
		log.Fatalf("Failed to select source directory: %v", err)
	}

	// Show file picker for destination directory
	destDir, err := dialog.Directory().Title("Select Destination Directory").Browse()
	if err != nil {
		log.Fatalf("Failed to select destination directory: %v", err)
	}

	// Validate source directory
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		log.Fatalf("Source folder not found: %v", err)
	}

	// Validate destination directory
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		log.Fatalf("Destination folder not found: %v", err)
	}

	sourceDrive, err := types.NewDrive(srcDir)
	if err != nil {
		log.Fatalf("Error accessing source drive: %v", err)
	}

	destDrive, err := types.NewDrive(destDir)
	if err != nil {
		log.Fatalf("Error accessing destination drive: %v", err)
	}

	fileTransfer := types.NewFileTransfer(
		*sourceDrive,
		*destDrive,
		true,
		func(src, dest string) {
			fmt.Printf("Transferred: %s → %s\n", src, dest)
		},
	)

	if err := fileTransfer.Transfer(); err != nil {
		log.Fatalf("Transfer error: %v", err)
	}

	fileOrganizer := types.NewFileOrganizer(
		*destDrive,
		func(src, dest string) {
			fmt.Printf("Organized: %s → %s\n", src, dest)
		},
	)

	if err := fileOrganizer.Organize(); err != nil {
		log.Fatalf("Organization error: %v", err)
	}

	fmt.Println("✅ Transfer Complete!")
}
