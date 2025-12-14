package mocks

import (
	"time"
	types "installer.malang/internal/types"
)

func CollectDisks() []types.Disk {
	time.Sleep(500 * time.Millisecond) // Simulate disk scanning
	return []types.Disk{
		{
			Name:       "sda",
			Size:       1000000000000,
			SizeInGb:   "1000G",
			Type:       "disk",
			MountPoint: "",
		},
		{
			Name:       "nvme0n1", 
			Size:       500000000000,
			SizeInGb:   "500G",
			Type:       "disk",
			MountPoint: "",
		},
		{
			Name:       "loop0",
			Size:       10737418240,
			SizeInGb:   "10G",
			Type:       "loop",
			MountPoint: "",
		},
	}
}