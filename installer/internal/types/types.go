package types

type Disk struct {
    Name             string
    Size             uint64
    SizeInGb     string
    Type             string
    MountPoint string
}

type ProgressUpdate struct {
    Message string
    Step        int
    Success bool
}

type InstallPackageStream struct {
    Line     string
    Source string
}

type ConfigureStream struct {
    Line string
}

type WiFiNetwork struct {
    SSID         string
    Security string
    Signal     int
    InUse        bool
    Hidden     bool

    BSSID         string
    Channel     int
    Frequency int

    RequiresPSK bool
}

type NetworkFailureMsg bool

type NetworkConnectedMsg bool

type WiFiNetworksMsg []WiFiNetwork

type WifiNetworkError string

type SelectedDiskMsg Disk

type StartPartitioningMsg bool

type PartitionConfigMsg [3]string

type PartitionError string

type InstallCompleteMsg [2]string
