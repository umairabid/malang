package mocks

import (
    types "installer.malang/internal/types"
    "time"
)

func PartitionDisk(disk types.Disk, percentages [3]int) ([3]string, error) {
    time.Sleep(2 * time.Second)

    if disk.Type == "loop" {
        return [3]string{
            "/dev/" + disk.Name + "p1",
            "/dev/" + disk.Name + "p2",
            "/dev/" + disk.Name + "p3",
        }, nil
    }

    return [3]string{
        "/dev/" + disk.Name + "1",
        "/dev/" + disk.Name + "2",
        "/dev/" + disk.Name + "3",
    }, nil
}
