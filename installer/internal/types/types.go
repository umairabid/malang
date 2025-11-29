package types

type Disk struct {
  Name       string
  Size       uint64
  SizeInGb   string
  Type       string
  MountPoint string
}

type SelectedDiskMsg Disk

type PartitionConfigMsg [3]string

type InstallCompleteMsg [2]string 

