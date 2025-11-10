package main

import (
	"fmt"
	//"os/exec"
	//"strings"
)

const (
	efiPercent  = 20
	swapPercent = 10
)

func getDriveNames(name string) [3]string {
	efiName := fmt.Sprintf("%sp1", name)
	swapName := fmt.Sprintf("%sp2", name)
	rootName := fmt.Sprintf("%sp3", name)

	driveNames := [...]string{efiName, swapName, rootName}
	return driveNames
}

func createScheme(name string, size uint64, driveNames [3]string) string {
	sizeInMb := size / 1024 / 1024

	efiSize := sizeInMb * efiPercent / 100
	swapSize := sizeInMb * swapPercent / 100
	rootSize := sizeInMb - (efiSize + swapSize)

	scheme := fmt.Sprintf(`label: gpt
device: %s

%s : start=, size=%dM, type=uefi
%s : start=, size=%dM, type=linux-swap
%s : start=, size=%dM, type=linux
`, name, driveNames[0], efiSize, driveNames[1], swapSize, driveNames[2], rootSize)
	fmt.Printf("--- Partition Scheme ---\n%s\n", scheme)
	return scheme
}

func partitionDisk(disk Disk) [3]string {
	diskName := disk.Name
	diskSize := disk.Size
	diskPath := "/dev/" + diskName

	driveNames := getDriveNames(diskPath)
	scheme := createScheme(diskPath, diskSize, driveNames)
	commands := []Command{
		Command{Args: []string{"wipefs", "-af", diskPath}},
		Command{Args: []string{"partprobe", diskPath}},
		Command{Args: []string{"sfdisk", "-f", diskPath}, Stdin: &scheme},
		Command{Args: []string{"partprobe", diskPath}},
	}
	runCommands(commands)

	fmt.Printf("✅ Partitioning succeeded for\n")

	return driveNames
}
