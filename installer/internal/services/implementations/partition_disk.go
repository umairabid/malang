package implementations

import (
	"fmt"
	"installer.malang/internal/utils"

  types "installer.malang/internal/types"
)

func sizes(totalSize uint64, percentages [3]int) [3]uint64 {
  sizeInMb := totalSize / 1024 / 1024

  bootSize := sizeInMb * uint64(percentages[0]) / 100
  swapSize := sizeInMb * uint64(percentages[1]) / 100
  rootSize := sizeInMb * uint64(percentages[2]) / 100

  return [3]uint64{bootSize, swapSize, rootSize}
}

func createScheme(name string, sizes [3]uint64) string {
	scheme := fmt.Sprintf(`label: gpt
device: %s

1 : start=, size=%dM, type=uefi
2 : start=, size=%dM, type=linux-swap
3 : start=, size=%dM, type=linux
`, name, sizes[0], sizes[1], sizes[2])
	return scheme
}

func PartitionDisk(disk types.Disk, percentages [3]int) ([3]string, error) {
  diskName := disk.Name
  diskSize := disk.Size

	diskPath := "/dev/" + diskName
  sizes := sizes(diskSize, percentages)

	resetCommands := []utils.Command{
		{Args: []string{"swapoff", "-a"}},
		{Args: []string{"umount", "-f", diskPath + "*"}},
		{Args: []string{"wipefs", "-af", diskPath}},
	}
	utils.RunCommands(resetCommands)

	scheme := createScheme(diskPath, sizes)
	commands := []utils.Command{
		{Args: []string{"sfdisk", "-f", diskPath}, Stdin: &scheme},
		{Args: []string{"partprobe", diskPath}}, // Add new partitions to kernel
	}
  err := utils.RunCommands(commands)
 
  if err != nil {
    return [3]string{}, err
  }

	driveNames := utils.FetchPartitions(diskName)
	return [3]string{driveNames[0], driveNames[1], driveNames[2]}, nil
}
