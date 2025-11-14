package main

import (
	"fmt"
)

const (
	efiPercent  = 20
	swapPercent = 10
)

func createScheme(name string, size uint64) string {
	sizeInMb := size / 1024 / 1024

	efiSize := sizeInMb * efiPercent / 100
	swapSize := sizeInMb * swapPercent / 100
	rootSize := sizeInMb - (efiSize + swapSize)

	scheme := fmt.Sprintf(`label: gpt
device: %s

1 : start=, size=%dM, type=uefi
2 : start=, size=%dM, type=linux-swap
3 : start=, size=%dM, type=linux
`, name, efiSize, swapSize, rootSize)
	fmt.Printf("--- Partition Scheme ---\n%s\n", scheme)
	return scheme
}

func partitionDisk(disk Disk) [3]string {
	diskName := disk.Name
	diskSize := disk.Size
	diskPath := "/dev/" + diskName

	scheme := createScheme(diskPath, diskSize)
	commands := []Command{
		Command{Args: []string{"wipefs", "-af", diskPath}},
		Command{Args: []string{"partprobe", diskPath}},
		Command{Args: []string{"sfdisk", "-f", diskPath}, Stdin: &scheme},
		Command{Args: []string{"partprobe", diskPath}},
	}
	runCommands(commands)

	fmt.Printf("✅ Partitioning succeeded for\n")

	driveNames := getPartitions(diskName)
	return [3]string{driveNames[0], driveNames[1], driveNames[2]}
}
