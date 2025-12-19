package services

import (
	"os"
	types "installer.malang/internal/types"
	"installer.malang/internal/services/implementations"
	"installer.malang/internal/services/mocks"
)

func Install(
	disks [3]string,
	progressChan chan<- types.ProgressUpdate,
	streamChan chan<- types.InstallPackageStream,
) ([2]string, error) {
	if os.Getenv("IS_MOCKING") == "true" {
		return mocks.Install(disks, progressChan, streamChan)
	}
	return implementations.Install(disks, progressChan, streamChan)
}
