package main

import (
	"fmt"
  //"installer.malang/internal/services"
  "installer.malang/internal/ui/app"
)

func main() {
  app.App()
	//disk := services.DiskForInstallation()
	//fmt.Printf("Selected disk for installation: %+v\n", disk)

	//driveNames := services.PartitionDisk(disk)
	//fmt.Printf("Created partitions: %v\n", driveNames)

	//mountPoints := services.Install(driveNames)
	fmt.Println("Installed Archlinux on disk.")

	//services.ConfigureSystem(mountPoints)
}
