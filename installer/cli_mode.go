package main

import (
    "fmt"
    services "installer.malang/internal/services"
    types "installer.malang/internal/types"
)

func RunCliMode() {
	networks, err := services.GetWiFiNetworks()
	fmt.Println(networks, err)

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
    if err != nil {
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

    mountPoints, err := services.Install(driveNames, progressChan, streamChan)
    fmt.Printf("Mount points: %v\n", mountPoints)
    if err != nil {
        fmt.Printf("Installation failed: %v\n", err)
        return
    }

    configureStream := make(chan types.ConfigureStream)
    go func() {
        for stream := range configureStream {
            fmt.Printf("Configuring system: %s\n", stream.Line)
        }
    }()
    services.ConfigureSystem(mountPoints, configureStream)
    fmt.Println("Installed Archlinux on disk.")
}
