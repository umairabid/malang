package services

import (
    "installer.malang/internal/services/implementations"
    "installer.malang/internal/services/mocks"
    "installer.malang/internal/types"
    "os"
)

func ConfigureSystem(mountPoints [2]string, progressChan chan types.ConfigureStream) {
    if os.Getenv("IS_MOCKING") == "true" {
        mocks.ConfigureSystem(mountPoints, progressChan)
        return
    }
    implementations.ConfigureSystem(mountPoints, progressChan)
}

func CreateUser(userConfig types.UserConfig, progressChan chan types.ConfigureStream) {
    if os.Getenv("IS_MOCKING") == "true" {
        mocks.CreateUser(userConfig, progressChan)
        return
    }
    implementations.CreateUser(userConfig, progressChan)
}
