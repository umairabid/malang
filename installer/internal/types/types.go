package types

type Disk struct {
  Name       string
  Size       uint64
  SizeInGb   string
  Type       string
  MountPoint string
}

type SelectedDiskMsg Disk

