package implementations

import (
    "bufio"
    "installer.malang/internal/utils"
    "io"
    "os"
    "os/exec"

    types "installer.malang/internal/types"
)

const RootMountPoint = "/mnt"
const BootDir = "boot/efi"
const BootMountPoint = RootMountPoint + "/" + BootDir

func preInstallSetup(disks [3]string) error {
    efiDisk := disks[0]
    swapDisk := disks[1]
    rootDisk := disks[2]

    prepareCommands := []utils.Command{
        {Args: []string{"mkfs.fat", "-F32", efiDisk}},
        {Args: []string{"mkfs.ext4", rootDisk}},
        {Args: []string{"mkswap", swapDisk}},
        {Args: []string{"swapon", swapDisk}},
        {Args: []string{"mount", rootDisk, RootMountPoint}},
        {Args: []string{"mkdir", "-p", BootMountPoint}},
        {Args: []string{"mount", efiDisk, BootMountPoint}},
    }

    return utils.RunCommands(prepareCommands)
}

func postInstallSetup() error {
    prepareCommands := []utils.Command{
        {Args: []string{"genfstab", "-U", "/mnt", ">>", RootMountPoint + "/etc/fstab"}},
        {Args: []string{"mount", "--types", "proc", "/proc", "/mnt/proc"}},
        {Args: []string{"mount", "--rbind", "/sys", "/mnt/sys"}},
        {Args: []string{"mount", "--rbind", "/dev", "/mnt/dev"}},
        {Args: []string{"mount", "--rbind", "/run", "/mnt/run"}},
    }
    return utils.RunCommands(prepareCommands)
}

func configureGreetd() error {
    configDir := RootMountPoint + "/etc/greetd"
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return err
    }

	greetdConfigPath := "/etc/greetd/config.toml"
    configPath := configDir + "/config.toml"
    return utils.RunCommands([]utils.Command{
		{Args: []string{"cp", greetdConfigPath, configPath}},
	})
}

func enableServices() error {
    // Configure greetd first
    if err := configureGreetd(); err != nil {
        return err
    }

    serviceCommands := []utils.Command{
        // Enable NetworkManager
        {Args: []string{"arch-chroot", RootMountPoint, "systemctl", "enable", "NetworkManager.service"}},
        // Enable greetd (display manager)
        {Args: []string{"arch-chroot", RootMountPoint, "systemctl", "enable", "greetd.service"}},
    }
    return utils.RunCommands(serviceCommands)
}

func startStream(pipe io.ReadCloser, streamName string, streamChan chan<- types.InstallPackageStream) {
    go func() {
        defer pipe.Close()
        scanner := bufio.NewScanner(pipe)
        for scanner.Scan() {
            streamChan <- types.InstallPackageStream{
                Line:     scanner.Text(),
                Source: streamName,
            }
        }
    }()
}

func InstallPackages(streamChan chan<- types.InstallPackageStream) error {
    cmd := exec.Command("pacstrap", "/mnt", "base", "linux", "linux-firmware", "vim", "networkmanager", "efibootmgr", "grub",
	"sway", "swaybg", "swaylock", "swayidle", "greetd", "greetd-tuigreet", "wofi", "waybar", "foot", "wmenu", "firefox", "sudo")
    stderr, _ := cmd.StderrPipe()
    stdout, _ := cmd.StdoutPipe()

    startStream(stderr, "stderr", streamChan)
    startStream(stdout, "stdout", streamChan)

    if err := cmd.Start(); err != nil {
        return err
    }

    return cmd.Wait()
}

func emitPackageInstallProgress(err error, progressChan chan<- types.ProgressUpdate) {
    if err != nil {
        progressChan <- types.ProgressUpdate{
            Message: "Failed to prepare disks: " + err.Error(),
            Step:        1,
            Success: false,
        }
    } else {
        progressChan <- types.ProgressUpdate{
            Message: "Installing packages.",
            Step:        2,
            Success: true,
        }
    }
}

func Install(
    disks [3]string,
    progressChan chan<- types.ProgressUpdate,
    streamChan chan<- types.InstallPackageStream,
) ([2]string, error) {
    progressChan <- types.ProgressUpdate{
        Message: "Starting installation process.",
        Step:        1,
        Success: true,
    }
    err := preInstallSetup(disks)
    emitPackageInstallProgress(err, progressChan)
    if err != nil {
        return [2]string{}, err
    }

    err = InstallPackages(streamChan)
    if err != nil {
        progressChan <- types.ProgressUpdate{
            Message: "Package installation failed: " + err.Error(),
            Step:        3,
            Success: false,
        }
        return [2]string{}, err
    }

    err = postInstallSetup()
    if err != nil {
        progressChan <- types.ProgressUpdate{
            Message: "Post-install setup failed: " + err.Error(),
            Step:        3,
            Success: false,
        }
        return [2]string{}, err
    }

    progressChan <- types.ProgressUpdate{
        Message: "Enabling system services.",
        Step:        3,
        Success: true,
    }

    err = enableServices()
    if err != nil {
        progressChan <- types.ProgressUpdate{
            Message: "Failed to enable services: " + err.Error(),
            Step:        4,
            Success: false,
        }
        return [2]string{}, err
    }

    progressChan <- types.ProgressUpdate{
        Message: "System installation completed successfully.",
        Step:        4,
        Success: true,
    }

    return [2]string{RootMountPoint, BootDir}, nil
}
