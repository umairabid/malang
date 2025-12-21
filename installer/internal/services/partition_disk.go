package services

import (
    "installer.malang/internal/services/implementations"
    "installer.malang/internal/services/mocks"
    types "installer.malang/internal/types"
    "os"
)

func PartitionDisk(disk types.Disk, percentages [3]int) ([3]string, error) {
    if os.Getenv("IS_MOCKING") == "true" {
        return mocks.PartitionDisk(disk, percentages)
    }
    return implementations.PartitionDisk(disk, percentages)
}
