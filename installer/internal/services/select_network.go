package services

import (
    "installer.malang/internal/services/mocks"
	  "installer.malang/internal/services/implementations"
    "os"

    types "installer.malang/internal/types"
)

func CheckNetworkConnection() bool {
    if os.Getenv("IS_MOCKING") == "true" {
        return mocks.CheckNetworkConnection()
    }
    return implementations.CheckNetworkConnection()
}

func GetWiFiNetworks() ([]types.WiFiNetwork, error) {
    if os.Getenv("IS_MOCKING") == "true" {
        return mocks.GetWiFiNetworks()
    }
    return implementations.GetWiFiNetworks()
}

func ConnectWithWiFi(ssid string, password string) error {
    if os.Getenv("IS_MOCKING") == "true" {
        return mocks.ConnectWithWiFi(ssid, password)
    }
    return implementations.ConnectWithWiFi(ssid, password)
}
