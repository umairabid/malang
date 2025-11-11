package main

func install(disks [3]string) [2]string {
	efiDisk := disks[0]
	swapDisk := disks[1]
	rootDisk := disks[2]

	rootMountPoint := "/mnt"
	bootDir := "boot/efi"
	bootMountPoint := rootMountPoint + "/" + bootDir

	prepareCommands := []Command{
		Command{Args: []string{"mkfs.fat", "-F32", efiDisk}},
		Command{Args: []string{"mkfs.ext4", rootDisk}},
		Command{Args: []string{"mkswap", swapDisk}},
		Command{Args: []string{"swapon", swapDisk}},
		Command{Args: []string{"mount", rootDisk, rootMountPoint}},
		Command{Args: []string{"mkdir", "-p", bootMountPoint}},
		Command{Args: []string{"mount", efiDisk, bootMountPoint}},
		Command{Args: []string{"pacstrap", "/mnt", "base", "linux", "linux-firmware", "vim", "networkmanager", "efibootmgr", "grub"}},
		Command{Args: []string{"genfstab", "-U", "/mnt", ">>", rootMountPoint + "/etc/fstab"}},
		Command{Args: []string{"mount", "--types", "proc", "/proc", "/mnt/proc"}},
		Command{Args: []string{"mount", "--rbind", "/sys", "/mnt/sys"}},
		Command{Args: []string{"mount", "--rbind", "/dev", "/mnt/dev"}},
		Command{Args: []string{"mount", "--rbind", "/run", "/mnt/run"}},
	}

	runCommands(prepareCommands)

	return [2]string{rootMountPoint, bootDir}
}
