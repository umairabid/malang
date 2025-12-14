package types

type Disk struct {
  Name       string
  Size       uint64
  SizeInGb   string
  Type       string
  MountPoint string
}

type ProgressUpdate struct {
	Message string
  Step    int
  Success bool
}

type InstallPackageStream struct {
  Line string
  Source string
}

type NetworkStatusMsg bool

type SelectedDiskMsg Disk

type StartPartitioningMsg bool

type PartitionConfigMsg [3]string

type PartitionError string

type InstallCompleteMsg [2]string 

