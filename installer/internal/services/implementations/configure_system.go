package implementations

import (
	"fmt"
	"os"
	"syscall"
  "installer.malang/internal/utils"
)

func ConfigureSystem(mountPoints [2]string) {
	rootMountPoint := mountPoints[0]
	bootDir := mountPoints[1]

	if err := syscall.Chroot(rootMountPoint); err != nil {
		panic("Failed to chroot:" + err.Error())
	}

	if err := os.Chdir("/"); err != nil {
		panic("Failed to chdir: %v" + err.Error())
	}

	fmt.Println("=== Installing GRUB ===")
	commands := []utils.Command{
		{Args: []string{"echo", "\"KEYMAP=us\" > /etc/vconsole.conf"}},
		{Args: []string{"grub-install", "--target=x86_64-efi", "--efi-directory=" + bootDir, "--bootloader-id=GRUB"}},
		{Args: []string{"grub-mkconfig", "-o", "/boot/grub/grub.cfg"}},
	}
	utils.RunCommands(commands)
}
