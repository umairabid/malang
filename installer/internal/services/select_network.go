package services

import (
	"os"
	"installer.malang/internal/services/mocks"
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
