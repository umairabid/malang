package mocks

import (
    "installer.malang/internal/types"
    "time"
)

func ConfigureSystem(mountPoints [2]string, progressChan chan types.ConfigureStream) {
    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "System configuration started."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Setting up locale..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Configuring network settings..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Installing bootloader..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Finalizing configuration..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "System configuration complete."}
}
