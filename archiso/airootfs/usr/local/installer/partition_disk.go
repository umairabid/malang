package main

import (
	"flag"
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

func createScheme(name string, size uint64) string {
  sizeInMb := size / 1024 / 1024

  efiSize := sizeInMb * efiPercent / 100
  swapSize := sizeInMb * swapPercent / 100
  rootSize := sizeInMb - (efiSize + swapSize)

  driveNames := getDriveNames(name)
  scheme := fmt.Sprintf(`label: gpt
device: %s

%s1 : start=, size=%dM, type=uefi
%s2 : start=, size=%dM, type=linux-swap
%s3 : start=, size=%dM, type=linux
`, name, driveNames[0], efiSize, driveNames[1], swapSize, driveNames[2], rootSize)

  return scheme
}

func main() {
	var diskName string
  var diskSize uint64
	flag.StringVar(&diskName, "diskName", "", "Name of disk to partition.")
  flag.Uint64Var(&diskSize, "diskSize", 0, "Size of selected disk")
	flag.Parse()

	diskPath := "/dev/" + diskName
	cmd := exec.Command("sfdisk", "-f", diskPath)

  scheme := createScheme(diskPath, diskSize)
  cmd.Stdin = strings.NewReader(scheme)
  output, err := cmd.CombinedOutput()

  if err != nil {
		fmt.Printf("--- sfdisk Output ---\n%s\n", string(output))
		fmt.Printf("Partitioning failed for %s. Error: %v", diskPath, err)
    os.Exit(1)
	}
  
  driveNames := getDriveNames(diskPath)
  fmt.Print(strings.Join(driveNames[:], "\n"))
}
