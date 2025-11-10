package main


func install(disks [3]string) [2]string {
  efiDisk := disks[0]
  swapDisk := disks[1]
  rootDisk := disks[2]

  rootMountPoint := "/mnt"
  bootMountPoint := "/mnt/boot/efi"

	prepareCommands := []Command{
    Command{Args: []string{"mkfs.fat", "-F32", efiDisk}},
    Command{Args: []string{"mkfs.ext4", rootDisk}},
    Command{Args: []string{"mkswap", swapDisk}},
    Command{Args: []string{"swapon", swapDisk}},
    Command{Args: []string{"mount", rootDisk, rootMountPoint}},
    Command{Args: []string{"mkdir", "-p", bootMountPoint}},
    Command{Args: []string{"mount", efiDisk, bootMountPoint}},
    Command{Args: []string{"pacstrap", "/mnt", "base", "linux", "linux-firmware", "vim", "networkmanager"}},
    Command{Args: []string{"genfstab", "-U", "/mnt", ">>", rootMountPoint + "/etc/fstab"}},
	}

  runCommands(prepareCommands)

  return [2]string {rootMountPoint, bootMountPoint}
}
