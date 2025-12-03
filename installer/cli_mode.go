package main

import (
  "fmt"
  services "installer.malang/internal/services"
  types "installer.malang/internal/types"
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
	driveNames, err := services.PartitionDisk(disk, percentages)
  if  err != nil {
    fmt.Printf("Failed to partition disk: %v\n", err)
    return
  }
	fmt.Printf("Created partitions: %v\n", driveNames)

  progressChan := make(chan types.ProgressUpdate, 10)
  streamChan := make(chan types.InstallPackageStream, 10)
  go func() {
    for update := range progressChan {
      fmt.Printf("Progress: Step %d - %s (Success: %v)\n", update.Step, update.Message, update.Success)
    }
  }()

  go func() {
    for stream := range streamChan {
      fmt.Printf("Installing package: %s - %s\n", stream.Line, stream.Source)
    }
  }()

	mountPoints := services.Install(driveNames, progressChan, streamChan)
  fmt.Printf("Mount points: %v\n", mountPoints)
	services.ConfigureSystem(mountPoints)
	fmt.Println("Installed Archlinux on disk.")
}
