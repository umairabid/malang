package services

import (
	"installer.malang/internal/utils"
)

func Install(disks [3]string) [2]string {
	efiDisk := disks[0]
	swapDisk := disks[1]
	rootDisk := disks[2]

	rootMountPoint := "/mnt"
	bootDir := "boot/efi"
	bootMountPoint := rootMountPoint + "/" + bootDir

	prepareCommands := []utils.Command{
		{Args: []string{"mkfs.fat", "-F32", efiDisk}},
		{Args: []string{"mkfs.ext4", rootDisk}},
		{Args: []string{"mkswap", swapDisk}},
		{Args: []string{"swapon", swapDisk}},
		{Args: []string{"mount", rootDisk, rootMountPoint}},
		{Args: []string{"mkdir", "-p", bootMountPoint}},
		{Args: []string{"mount", efiDisk, bootMountPoint}},
		{Args: []string{"pacstrap", "/mnt", "base", "linux", "linux-firmware", "vim", "networkmanager", "efibootmgr", "grub"}},
		{Args: []string{"genfstab", "-U", "/mnt", ">>", rootMountPoint + "/etc/fstab"}},
		{Args: []string{"mount", "--types", "proc", "/proc", "/mnt/proc"}},
		{Args: []string{"mount", "--rbind", "/sys", "/mnt/sys"}},
		{Args: []string{"mount", "--rbind", "/dev", "/mnt/dev"}},
		{Args: []string{"mount", "--rbind", "/run", "/mnt/run"}},
	}

	utils.RunCommands(prepareCommands)

	return [2]string{rootMountPoint, bootDir}
}
