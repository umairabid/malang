package main

import (
	"fmt"
)

func main() {
	disk := diskForInstallation()
	driveNames := partitionDisk(disk)
	install(driveNames)
	fmt.Printf("Selected disk for installation: %+v\n", disk)
	fmt.Printf("Created partitions: %v\n", driveNames)
	fmt.Println("This is a placeholder for the main package.")
}
