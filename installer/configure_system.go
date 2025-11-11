package main

import (
	"fmt"
	"os"
	"syscall"
)

func configureSystem(mountPoints [2]string) {
	rootMountPoint := mountPoints[0]
	bootDir := mountPoints[1]

	if err := syscall.Chroot(rootMountPoint); err != nil {
		panic("Failed to chroot:" + err.Error())
	}

	if err := os.Chdir("/"); err != nil {
		panic("Failed to chdir: %v" + err.Error())
	}

	fmt.Println("=== Installing GRUB ===")
	commands := []Command{
		Command{Args: []string{"echo", "\"KEYMAP=us\" > /etc/vconsole.conf"}},
		Command{Args: []string{"grub-install", "--target=x86_64-efi", "--efi-directory=" + bootDir, "--bootloader-id=GRUB"}},
		Command{Args: []string{"grub-mkconfig", "-o", "/boot/grub/grub.cfg"}},
	}
	runCommands(commands)
}
