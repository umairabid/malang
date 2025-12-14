package services

import (
	"os"
	"installer.malang/internal/services/implementations"
	"installer.malang/internal/services/mocks"
)

func ConfigureSystem(mountPoints [2]string) {
	if os.Getenv("IS_MOCKING") == "true" {
		mocks.ConfigureSystem(mountPoints)
		return
	}
	implementations.ConfigureSystem(mountPoints)
}