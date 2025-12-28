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

    progressChan <- types.ConfigureStream{Line: "Configuration complete"}
}

func CreateUser(userConfig types.UserConfig, progressChan chan types.ConfigureStream) {
    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Creating user " + userConfig.Username + "..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Setting up home directory..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "Adding user to wheel group..."}

    time.Sleep(1 * time.Second)

    progressChan <- types.ConfigureStream{Line: "User created successfully"}
}
