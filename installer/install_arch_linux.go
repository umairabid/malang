package main

func install(disks [3]string) {
	efiDisk := disks[0]
	swapDisk := disks[1]
	rootDisk := disks[2]

	// format and setup swap
	prepareCommands := []Command{
		Command{Args: []string{"mkfs.fat", "-F32", efiDisk}},
		Command{Args: []string{"mkfs.ext4", rootDisk}},
		Command{Args: []string{"mkswap", swapDisk}},
		Command{Args: []string{"swapon", swapDisk}},
	}

	runCommands(prepareCommands)
}
