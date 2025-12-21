package services

import (
    "installer.malang/internal/services/implementations"
    "installer.malang/internal/services/mocks"
    "os"

    types "installer.malang/internal/types"
)

func CollectDisks() []types.Disk {
    if os.Getenv("IS_MOCKING") == "true" {
        return mocks.CollectDisks()
    }
    return implementations.CollectDisks()
}
