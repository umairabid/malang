package main

import (
  "fmt"
  "os/exec"
  "strings"
  "os"
)

const (
  efiPercent = 20
  swapPercent = 10
)

func getDriveNames(name string) [3]string {
  efiName := fmt.Sprintf("%s1", name)
  swapName := fmt.Sprintf("%s2", name)
  rootName := fmt.Sprintf("%s3", name)

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
	cmd := exec.Command("sfdisk", "-f", diskPath)

  driveNames := getDriveNames(diskPath)
  scheme := createScheme(diskPath, diskSize, driveNames)
  cmd.Stdin = strings.NewReader(scheme)
  output, err := cmd.CombinedOutput()

  if err != nil {
		fmt.Printf("--- sfdisk Output ---\n%s\n", string(output))
		fmt.Printf("Partitioning failed for %s. Error: %v", diskPath, err)
    os.Exit(1)
	}
  
  return driveNames
}
