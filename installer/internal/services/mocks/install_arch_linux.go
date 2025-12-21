package mocks

import (
    "fmt"
    "time"

    types "installer.malang/internal/types"
)

func Install(
    disks [3]string,
    progressChan chan<- types.ProgressUpdate,
    streamChan chan<- types.InstallPackageStream,
) ([2]string, error) {
    fmt.Println("Mock Install called with disks:", disks)
    progressChan <- types.ProgressUpdate{
        Message: "Starting mock installation process.",
        Step:        1,
        Success: true,
    }

    packages := []string{
        "base", "linux", "linux-firmware", "vim",
        "networkmanager", "efibootmgr", "grub",
    }

    progressChan <- types.ProgressUpdate{
        Message: "Installing packages (mock).",
        Step:        2,
        Success: true,
    }

    for _, pkg := range packages {
        streamChan <- types.InstallPackageStream{
            Line:     "Installing package: " + pkg,
            Source: "stdout",
        }
        time.Sleep(200 * time.Millisecond)
    }

    progressChan <- types.ProgressUpdate{
        Message: "Mock system installation completed successfully.",
        Step:        3,
        Success: true,
    }

    return [2]string{"/mnt", "boot/efi"}, nil
}
