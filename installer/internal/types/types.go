package types

type Disk struct {
  Name       string
  Size       uint64
  SizeInGb   string
  Type       string
  MountPoint string
}

type SelectedDiskMsg Disk

type PartitionConfigMsg struct {
  Disk        Disk
  Percentages [3]int // [boot, swap, root]
}

