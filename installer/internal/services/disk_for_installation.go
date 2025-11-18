package services

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Disk struct {
	Name       string
	Size       uint64
	SizeInGb   string
	Type       string
	MountPoint string
}

func listDevices() []string {
	cmd := exec.Command("lsblk", "-b", "-d", "-o", "NAME,SIZE,TYPE,MOUNTPOINTS")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return strings.Split(string(out), "\n")
}

func parseDisks(lines []string) []Disk {
	var disks []Disk
	fmt.Println(lines)
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		name := fields[0]
		size, _ := strconv.ParseUint(fields[1], 10, 64)
		sizeInGb := size / (1024 * 1024 * 1024)
		formattedSizeInGb := strconv.FormatUint(sizeInGb, 10) + "G"

		dtype := fields[2]
		var mountPoint string

		if 3 < len(fields) {
			mountPoint = fields[3]
		}

		if sizeInGb > 0 {
			disks = append(disks, Disk{Name: name, Size: size, Type: dtype, MountPoint: mountPoint, SizeInGb: formattedSizeInGb})
		}
	}

	return disks
}

func collectDisks() []Disk {
	lines := listDevices()
	return parseDisks(lines)
}

func DiskForInstallation() Disk {
	disks := collectDisks()
	if len(disks) == 0 {
		fmt.Println("No devices found")
		os.Exit(1)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .SizeInGb | red }} bytes, {{ .MountPoint }})",
		Inactive: "  {{ .Name | cyan }} ({{ .SizeInGb }} bytes, {{ .MountPoint }})",
		Selected: "\U0001F389 Selected: {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select a disk",
		Items:     disks,
		Templates: templates,
		Size:      5,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Println("No device selected")
		os.Exit(1)
	}

	return disks[i]
}
