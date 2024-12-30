package main

import (
	"ShutterSync/internal/types"
	"fmt"
	"log"
)

func main() {
	source := "/Volumes/Untitled/DCIM/100MSDCF"
	dest := "/Volumes/My Passport/Photos"

	sourceDrive, err := types.NewDrive(source)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %v", err, source)
	}

	destDrive, err := types.NewDrive(dest)
	if err != nil {
		log.Fatalf("Error: %v. Failed to create a drive at path: %v", err, dest)
	}

	fileTransfer := types.NewFileTransfer(
		*sourceDrive,
		*destDrive,
		true,
		true,
		func(src, dest string) {
			fmt.Printf("Transferred %s to %s\n", src, dest)
		},
	)

	fileTransfer.Transfer()

	organizer := types.NewFileOrganizer(*destDrive, func(src, dest string) {
		fmt.Printf("Transferred %s to %s\n", src, dest)
	})

	err = organizer.Organize()

	fmt.Println(err)

}
