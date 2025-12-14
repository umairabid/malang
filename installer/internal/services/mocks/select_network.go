package mocks

import (
	"time"
)

func SelectNetwork() (bool, error) {
	time.Sleep(2 * time.Second)

  return false, nil
}

func CheckNetworkConnection() bool {
  time.Sleep(2 * time.Second)

  return true
}
