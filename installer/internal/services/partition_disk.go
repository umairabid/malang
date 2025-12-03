package services

import (
	"fmt"
	"installer.malang/internal/utils"

  types "installer.malang/internal/types"
)

func createScheme(name string, size uint64, percentages [3]int) string {
	sizeInMb := size / 1024 / 1024

	bootSize := sizeInMb * uint64(percentages[0]) / 100
	swapSize := sizeInMb * uint64(percentages[1]) / 100
	rootSize := sizeInMb * uint64(percentages[2]) / 100

	scheme := fmt.Sprintf(`label: gpt
device: %s

1 : start=, size=%dM, type=uefi
2 : start=, size=%dM, type=linux-swap
3 : start=, size=%dM, type=linux
`, name, bootSize, swapSize, rootSize)
	return scheme
}

func PartitionDisk(disk types.Disk, percentages [3]int) ([3]string, error) {
  diskName := disk.Name
  diskSize := disk.Size

	diskPath := "/dev/" + diskName
	scheme := createScheme(diskPath, diskSize, percentages)
	commands := []utils.Command{
		{Args: []string{"wipefs", "-af", diskPath}},
		{Args: []string{"partprobe", diskPath}},
		{Args: []string{"sfdisk", "-f", diskPath}, Stdin: &scheme},
		{Args: []string{"partprobe", diskPath}},
	}
  err := utils.RunCommands(commands)
  if err != nil {
    return [3]string{}, err
  }

	driveNames := utils.FetchPartitions(diskName)
	return [3]string{driveNames[0], driveNames[1], driveNames[2]}, nil
}
