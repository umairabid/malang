package implementations

import (
    "installer.malang/internal/types"
    "installer.malang/internal/utils"
    "os"
    "syscall"
)

func ConfigureSystem(mountPoints [2]string, progressChan chan types.ConfigureStream) error {
    rootMountPoint := mountPoints[0]
    bootDir := mountPoints[1]

    progressChan <- types.ConfigureStream{Line: "Changing root to installed system..."}
    if err := syscall.Chroot(rootMountPoint); err != nil {
        return err
    }

    progressChan <- types.ConfigureStream{Line: "Changing working directory to root..."}
    if err := os.Chdir("/"); err != nil {
        return err
    }

    progressChan <- types.ConfigureStream{Line: "Installing bootloader..."}
    commands := []utils.Command{
        {Args: []string{"echo", "\"KEYMAP=us\" > /etc/vconsole.conf"}},
        {Args: []string{"grub-install", "--target=x86_64-efi", "--efi-directory=" + bootDir, "--bootloader-id=GRUB"}},
        {Args: []string{"grub-mkconfig", "-o", "/boot/grub/grub.cfg"}},
    }
    err := utils.RunCommands(commands)
    return err
}
