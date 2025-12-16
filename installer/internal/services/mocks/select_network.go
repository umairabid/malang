package mocks

import (
	"time"

  types "installer.malang/internal/types"
)

func SelectNetwork() (bool, error) {
	time.Sleep(2 * time.Second)

  return false, nil
}

func CheckNetworkConnection() bool {
  time.Sleep(2 * time.Second)

  return true
}

func GetWiFiNetworks() ([]types.WiFiNetwork, error) {
  time.Sleep(2 * time.Second)

  networks := []types.WiFiNetwork{
    {
      SSID:     "HomeWiFi",
      Security: "WPA2",
      Signal:   85,
      InUse:    false,
      Hidden:   false,
      BSSID:    "00:11:22:33:44:55",
      Channel:  6,
      Frequency: 2437,
      RequiresPSK: true,
    },
    {
      SSID:     "CafeHotspot",
      Security: "OPEN",
      Signal:   70,
      InUse:    false,
      Hidden:   false,
      BSSID:    "66:77:88:99:AA:BB",
      Channel:  11,
      Frequency: 2462,
      RequiresPSK: false,
    },
    {
      SSID:     "OfficeWiFi",
      Security: "WPA3",
      Signal:   60,
      InUse:    false,
      Hidden:   false,
      BSSID:    "CC:DD:EE:FF:00:11",
      Channel:  1,
      Frequency: 2412,
      RequiresPSK: true,
    },
    {
      SSID:     "GuestNetwork",
      Security: "WPA2 WPA3",
      Signal:   50,
      InUse:    false,
      Hidden:   false,
      BSSID:    "22:33:44:55:66:77",
      Channel:  36,
      Frequency: 5180,
      RequiresPSK: true,
    },
  }

  return networks, nil
}

func ConnectWithWiFi(ssid string, password string) error {
  time.Sleep(2 * time.Second)

  return nil
}
