package services

import (
	"bufio"
	"fmt"
	"installer.malang/internal/utils"
	"os/exec"
)

type ProgressUpdate struct {
	Current int
	Total   int
	Message string
}

const RootMountPoint = "/mnt"
const BootDir = "boot/efi"
const BootMountPoint = RootMountPoint + "/" + BootDir

func preInstallSetup(disks [3]string) {
	efiDisk := disks[0]
	rootDisk := disks[1]
	swapDisk := disks[2]

	prepareCommands := []utils.Command{
		{Args: []string{"mkfs.fat", "-F32", efiDisk}},
		{Args: []string{"mkfs.ext4", rootDisk}},
		{Args: []string{"mkswap", swapDisk}},
		{Args: []string{"swapon", swapDisk}},
		{Args: []string{"mount", rootDisk, RootMountPoint}},
		{Args: []string{"mkdir", "-p", BootMountPoint}},
		{Args: []string{"mount", efiDisk, BootMountPoint}},
	}

	utils.RunCommands(prepareCommands)
}

func postInstallSetup() {
	prepareCommands := []utils.Command{
		{Args: []string{"genfstab", "-U", "/mnt", ">>", RootMountPoint + "/etc/fstab"}},
		{Args: []string{"mount", "--types", "proc", "/proc", "/mnt/proc"}},
		{Args: []string{"mount", "--rbind", "/sys", "/mnt/sys"}},
		{Args: []string{"mount", "--rbind", "/dev", "/mnt/dev"}},
		{Args: []string{"mount", "--rbind", "/run", "/mnt/run"}},
	}
	utils.RunCommands(prepareCommands)
}

func InstallPackages() {
	cmd := exec.Command("pacstrap", "/mnt", "base", "linux", "linux-firmware", "vim", "networkmanager", "efibootmgr", "grub")
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	reader := bufio.NewReader(stderr)
	var line []byte
	var totalPackages int
	var currentPackage int

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		// Carriage return or newline indicates end of current output
		if b == '\r' || b == '\n' {
			if len(line) > 0 {
				lineStr := string(line)

				// Parse package progress: "(X/Y) installing/downloading..."
				var current, total int
				if n, _ := fmt.Sscanf(lineStr, "(%d/%d)", &current, &total); n == 2 {
					totalPackages = total
					currentPackage = current

					// Calculate and print overall percentage
					if totalPackages > 0 {
						overallPercent := (currentPackage * 100) / totalPackages
						fmt.Printf("\rProgress: %d%% - %s", overallPercent, lineStr)
					}
				} else {
					// Print other output normally
					fmt.Print("\r" + lineStr)
				}

				line = line[:0]
			}
		} else {
			line = append(line, b)
		}
	}

	fmt.Println("\nInstallation complete: 100%")
	cmd.Wait()
}

func Install(disks [3]string, progressChan chan<- ProgressUpdate) [2]string {
  preInstallSetup(disks)
  InstallPackages()
  postInstallSetup()

  /**
  if progressChan != nil {
    progressChan <- ProgressUpdate{
      Current: i + 1,
      Total:   total,
      Message: commandDescriptions[i],
    }
  }
  **/
  return [2]string{RootMountPoint, BootDir}
}
