package implementations

import (
	"fmt"
	"os"
	"syscall"
  "installer.malang/internal/utils"
  "installer.malang/internal/types"
)

func ConfigureSystem(mountPoints [2]string, progressChan chan types.ConfigureStream) error {
	rootMountPoint := mountPoints[0]
	bootDir := mountPoints[1]

	if err := syscall.Chroot(rootMountPoint); err != nil {
    return err
	}

	if err := os.Chdir("/"); err != nil {
    return err
	}

	fmt.Println("=== Installing GRUB ===")
	commands := []utils.Command{
		{Args: []string{"echo", "\"KEYMAP=us\" > /etc/vconsole.conf"}},
		{Args: []string{"grub-install", "--target=x86_64-efi", "--efi-directory=" + bootDir, "--bootloader-id=GRUB"}},
		{Args: []string{"grub-mkconfig", "-o", "/boot/grub/grub.cfg"}},
	}
  err := utils.RunCommands(commands)
  return err
}
