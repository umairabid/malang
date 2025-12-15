package services

import (
	"os"
	"installer.malang/internal/services/mocks"

  types "installer.malang/internal/types"
)

func SelectNetwork() (bool, error) {
	if os.Getenv("IS_MOCKING") == "true" {
		return mocks.SelectNetwork()
	}
	return false, nil 
}

func CheckNetworkConnection() bool {
  if os.Getenv("IS_MOCKING") == "true" {
    return mocks.CheckNetworkConnection()
  }
  return false
}

func GetWiFiNetworks() ([]types.WiFiNetwork, error) {
  if os.Getenv("IS_MOCKING") == "true" {
    return mocks.GetWiFiNetworks()
  }
  return []types.WiFiNetwork{}, nil
}
