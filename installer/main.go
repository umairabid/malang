package main

import (
	"fmt"
)

func main() {
	disk := diskForInstallation()
	fmt.Printf("Selected disk for installation: %+v\n", disk)

	driveNames := partitionDisk(disk)
	fmt.Printf("Created partitions: %v\n", driveNames)

	mountPoints := install(driveNames)
	fmt.Println("Installed Archlinux on disk.")

	configureSystem(mountPoints)
}
