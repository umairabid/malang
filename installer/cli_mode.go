package main

import (
  "fmt"
  services "installer.malang/internal/services"
)

func RunCliMode() {
  disks := services.CollectDisks()
  fmt.Println("Available disks, select by index:")
  for i, disk := range disks {
    fmt.Printf("[%d]: %+v\n", i, disk)
  }
  
  var index int
  fmt.Print("Enter disk index: ")
  fmt.Scanln(&index)
  
  if index < 0 || index >= len(disks) {
    fmt.Println("Invalid index")
    return
  }


  disk := disks[index]
	fmt.Printf("Selected disk for installation: %+v\n", disk)
  
  percentages := [3]int{20, 20, 70}
	driveNames := services.PartitionDisk(disk, percentages)
	fmt.Printf("Created partitions: %v\n", driveNames)

  progressChan := make(chan services.ProgressUpdate, 10)
	mountPoints := services.Install(driveNames, progressChan)
  fmt.Printf("Mount points: %v\n", mountPoints)
	services.ConfigureSystem(mountPoints)
	fmt.Println("Installed Archlinux on disk.")
}
